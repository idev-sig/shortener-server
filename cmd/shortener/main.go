/*
Copyright © 2025 Jetsung Chan <i@jetsung.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"resty.dev/v3"

	"go.bdev.cn/shortener/internal/types"
)

// Config 配置
type Config struct {
	APIURL string `mapstructure:"url"`
	APIKEY string `mapstructure:"key"`
}

const (
	appName    = "shortener"
	configName = "config"
	configType = "toml"
	cfgDirName = "shortener"
	version    = "0.2.2"
)

var (
	cfg           Config
	configDir     string
	APIRequestURL = "/api"
	APIShortenURL = "/shortens"
	rootCmd       = &cobra.Command{
		Use:           appName,
		Short:         "Short URL management CLI tool",
		SilenceUsage:  true,  // 隐藏错误时自动显示的用法帮助
		SilenceErrors: false, // 隐藏Cobra默认的错误前缀
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := initConfig(); err != nil {
				return err
			}
			return nil
		},
	}
)

func init() {
	viper.SetEnvPrefix("SHORTENER")
	viper.AutomaticEnv()

	rootCmd.PersistentFlags().StringP("url", "u", "", "API URL")
	rootCmd.PersistentFlags().StringP("key", "k", "", "API KEY")

	_ = viper.BindPFlag("url", rootCmd.PersistentFlags().Lookup("url"))
	_ = viper.BindPFlag("key", rootCmd.PersistentFlags().Lookup("key"))

	rootCmd.AddCommand(newInitCmd())
	rootCmd.AddCommand(newEnvCmd())
	rootCmd.AddCommand(newVersionCmd())
	rootCmd.AddCommand(newShortenCreateCmd())
	rootCmd.AddCommand(newShortenDeleteCmd())
	rootCmd.AddCommand(newShortenUpdateCmd())
	rootCmd.AddCommand(newShortenGetCmd())
	rootCmd.AddCommand(newShortenListCmd())
}

func initConfig() error {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	configDir = filepath.Join(userConfigDir, cfgDirName)
	// log.Printf("configDir: %s", configDir)

	viper.SetConfigName(configName)
	viper.AddConfigPath(configDir)
	viper.SetConfigType(configType)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("failed to read config file: %w", err)
		}
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if cfg.APIKEY == "" {
		return nil
	}

	APIRequestURL = cfg.APIURL + APIRequestURL
	APIShortenURL = APIRequestURL + APIShortenURL

	return nil
}

// IsURL 判断是否为URL
func isURL(url string) bool {
	return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
}

func checkConfig() error {
	if !isURL(APIRequestURL) {
		return errors.New(`
  必须提供API地址，可用方式：
    1. 命令行参数: --url
    2. 环境变量: export SHORTENER_URL=your_url
    3. 配置文件: 在 ~/.shortener/config.toml 添加 url
	`)
	}
	return nil
}

func newInitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Initialize configuration",
		Run: func(cmd *cobra.Command, args []string) {
			if err := os.MkdirAll(configDir, 0o700); err != nil {
				fmt.Printf("Create config directory failed: %s\n%v\n", configDir, err)
				return
			}

			apiURL := cmd.Flags().Lookup("url").Value.String()
			apiKey := cmd.Flags().Lookup("key").Value.String()

			if apiURL != "" {
				viper.Set("url", apiURL)
			}
			if apiKey != "" {
				viper.Set("key", apiKey)
			}

			configFile := filepath.Join(configDir, configName+"."+configType)
			if err := viper.WriteConfigAs(configFile); err != nil {
				fmt.Printf("Write config failed: %s\n%v\n", configFile, err)
				return
			}

			fmt.Printf("Config file created: %s\n", configFile)
		},
	}
}

func newEnvCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "env",
		Aliases: []string{"e"},
		Short:   "Print environment variables",
		Example: `  shortener env`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("SHORTENER_URL: %s\n", viper.GetString("url"))
			fmt.Printf("SHORTENER_KEY: %s\n", viper.GetString("key"))
		},
	}
}

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "version",
		Aliases: []string{"v"},
		Short:   "Print version information",
		Example: `  shortener version`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("shortener: %s\n", version)
		},
	}
}

func newShortenCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create <origin_url>",
		Aliases: []string{"add"},
		Short:   "Create a short link",
		Args:    cobra.ExactArgs(1),
		Example: `  shortener create https://example.com/long/url
  shortener create https://example.com --code CUSTOM_CODE --desc "My special link"`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return checkConfig()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("origin URL is required")
			}

			originURL := args[0]
			if !isURL(originURL) {
				return fmt.Errorf("invalid origin URL: %s", originURL)
			}

			customCode, _ := cmd.Flags().GetString("code")
			description, _ := cmd.Flags().GetString("desc")

			req := struct {
				Code        string `json:"code,omitempty"`
				OriginalURL string `json:"original_url" binding:"required"`
				Describe    string `json:"describe,omitempty"`
			}{
				Code:        customCode,
				OriginalURL: originURL,
				Describe:    description,
			}

			client := resty.New()
			defer client.Close()

			var response types.ResShorten
			var resErr types.ResErr

			res, err := client.R().
				SetHeader("X-API-KEY", cfg.APIKEY).
				SetContentType("application/json").
				SetBody(req).
				SetResult(&response).
				SetError(&resErr).
				Post(APIShortenURL)
			if err != nil {
				return fmt.Errorf("failed to create short URL: \n  %w", err)
			}

			if res.StatusCode() != http.StatusCreated {
				return fmt.Errorf("failed to create short URL: \n  status code: %d \n      errcode: %d \n      errinfo: %s",
					res.StatusCode(),
					resErr.ErrCode,
					resErr.ErrInfo)
			}

			if !isURL(response.OriginalURL) {
				return fmt.Errorf("invalid short url: %s", response.OriginalURL)
			}

			fmt.Printf("Created short Code: %s\n", response.Code)
			fmt.Printf("         Short URL: %s\n", response.ShortURL)
			fmt.Printf("      Original URL: %s\n", response.OriginalURL)
			fmt.Printf("       Description: %s\n", response.Describe)
			return nil
		},
	}

	cmd.Flags().StringP("code", "c", "", "Custom short code (optional)")
	cmd.Flags().StringP("desc", "d", "", "Link description (optional)")

	return cmd
}

func newShortenDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "delete <short_code>",
		Aliases: []string{"del"},
		Short:   "Delete a short code",
		Args:    cobra.ExactArgs(1),
		Example: `  shortener delete MySpecialCode`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return checkConfig()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("short code is required")
			}

			code := args[0]

			client := resty.New()
			defer client.Close()

			var resErr types.ResErr

			res, err := client.R().
				SetHeader("X-API-KEY", cfg.APIKEY).
				SetContentType("application/json").
				SetError(&resErr).
				Delete(APIShortenURL + "/" + code)
			if err != nil {
				return fmt.Errorf("failed to delete short URL: \n  %w", err)
			}

			if res.StatusCode() != 204 {
				return fmt.Errorf("failed to delete short URL: \n  status code: %d \n      errcode: %d \n      errinfo: %s",
					res.StatusCode(),
					resErr.ErrCode,
					resErr.ErrInfo)
			}

			fmt.Printf("Deleted short Code: %s\n", code)
			return nil
		},
	}
}

func newShortenUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update <short_code>",
		Short: "Update a short code",
		Args:  cobra.ExactArgs(1),
		Example: `  shortener update MySpecialCode --ourl https://example.com
  shortener update MySpecialCode --ourl https://example.com --desc "My special link"`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return checkConfig()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("short code is required")
			}

			code := args[0]

			originURL, _ := cmd.Flags().GetString("ourl")
			description, _ := cmd.Flags().GetString("desc")

			if originURL != "" && !isURL(originURL) {
				return fmt.Errorf("invalid origin URL: %s", originURL)
			}

			req := struct {
				OriginalURL string `json:"original_url,omitempty" binding:"omitempty,url"`
				Describe    string `json:"describe,omitempty"`
			}{
				OriginalURL: originURL,
				Describe:    description,
			}

			var response types.ResShorten
			var resErr types.ResErr

			client := resty.New()
			defer client.Close()

			res, err := client.R().
				SetHeader("X-API-KEY", cfg.APIKEY).
				SetContentType("application/json").
				SetBody(req).
				SetResult(&response).
				SetError(&resErr).
				Put(APIShortenURL + "/" + code)
			if err != nil {
				return fmt.Errorf("failed to update short code: \n  %w", err)
			}

			if res.StatusCode() != http.StatusOK {
				return fmt.Errorf("failed to update short URL: \n  status code: %d \n      errcode: %d \n      errinfo: %s",
					res.StatusCode(),
					resErr.ErrCode,
					resErr.ErrInfo)
			}

			if !isURL(response.OriginalURL) {
				return fmt.Errorf("invalid short code: %s", code)
			}

			fmt.Printf("Updated short Code: %s\n", response.Code)
			fmt.Printf("         Short URL: %s\n", response.ShortURL)
			fmt.Printf("      Original URL: %s\n", response.OriginalURL)
			fmt.Printf("       Description: %s\n", response.Describe)
			return nil
		},
	}

	cmd.Flags().StringP("ourl", "o", "", "Original URL (optional)")
	cmd.Flags().StringP("desc", "d", "", "Link description (optional)")

	return cmd
}

func newShortenGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "get <short_code>",
		Aliases: []string{"g"},
		Short:   "Get a short link",
		Args:    cobra.ExactArgs(1),
		Example: `  shortener get MySpecialCode`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return checkConfig()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("short code is required")
			}

			code := args[0]

			client := resty.New()
			defer client.Close()

			var response types.ResShorten
			var resErr types.ResErr

			res, err := client.R().
				SetHeader("X-API-KEY", cfg.APIKEY).
				SetContentType("application/json").
				SetResult(&response).
				SetError(&resErr).
				Get(APIShortenURL + "/" + code)
			if err != nil {
				return fmt.Errorf("failed to get short URL: \n  %w", err)
			}

			if res.StatusCode() != http.StatusOK {
				return fmt.Errorf("failed to get short URL: \n  status code: %d \n      errcode: %d \n      errinfo: %s",
					res.StatusCode(),
					resErr.ErrCode,
					resErr.ErrInfo)
			}

			if !isURL(response.OriginalURL) {
				return fmt.Errorf("invalid short code: %s", code)
			}

			fmt.Printf("  Short Code: %s\n", response.Code)
			fmt.Printf("   Short URL: %s\n", response.ShortURL)
			fmt.Printf("Original URL: %s\n", response.OriginalURL)
			fmt.Printf(" Description: %s\n", response.Describe)
			return nil
		},
	}
}

func newShortenListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"l"},
		Short:   "List all short links",
		Example: `  shortener list`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return checkConfig()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			isAll, _ := cmd.Flags().GetBool("all")

			client := resty.New()
			defer client.Close()

			var allData []types.ResShorten
			var response types.ResSuccess[[]types.ResShorten]
			var resErr types.ResErr

			// 处理获取所有数据的逻辑
			if isAll {
				page := int64(1)
				pageSize := int64(100) // 每页获取最大允许数量

				for {
					query := url.Values{}
					query.Set("page", strconv.FormatInt(page, 10))
					query.Set("page_size", strconv.FormatInt(pageSize, 10))
					query.Set("sort_by", "created_at")
					query.Set("order", "asc")

					res, err := client.R().
						SetHeader("X-API-KEY", cfg.APIKEY).
						SetContentType("application/json").
						SetResult(&response).
						SetError(&resErr).
						Get(APIShortenURL + "?" + query.Encode())
					if err != nil {
						return fmt.Errorf("failed to list short URLs: \n  %w", err)
					}

					if res.StatusCode() != http.StatusOK {
						return fmt.Errorf("failed to list short URLs: \n  status code: %d \n      errcode: %d \n      errinfo: %s",
							res.StatusCode(),
							resErr.ErrCode,
							resErr.ErrInfo)
					}

					// 合并数据
					allData = append(allData, response.Data...)

					// 检查是否还有更多数据
					if page >= response.Meta.TotalPages {
						break
					}
					page++
				}

				// 替换response为完整数据
				response.Data = allData
				response.Meta.Page = 1
				response.Meta.CurrentCount = int64(len(allData))
				response.Meta.TotalItems = int64(len(allData))
				response.Meta.TotalPages = 1
				response.Meta.PageSize = int64(len(allData))
			} else {
				// 原有分页查询逻辑
				page, _ := cmd.Flags().GetInt64("page")
				pageSize, _ := cmd.Flags().GetInt64("psize")
				sortBy, _ := cmd.Flags().GetString("sort")
				order, _ := cmd.Flags().GetString("order")

				// 搜索
				code, _ := cmd.Flags().GetString("code")
				originalURL, _ := cmd.Flags().GetString("original_url")

				// 设置默认值
				if page == 0 {
					page = 1
				}
				if pageSize == 0 {
					pageSize = 10
				}
				if sortBy == "" {
					sortBy = "created_at"
				}
				if order == "" {
					order = "asc"
				}

				query := url.Values{}
				query.Set("page", strconv.FormatInt(page, 10))
				query.Set("page_size", strconv.FormatInt(pageSize, 10))
				query.Set("sort_by", sortBy)
				query.Set("order", order)

				if code != "" {
					query.Set("code", code)
				}
				if originalURL != "" {
					query.Set("original_url", originalURL)
				}

				res, err := client.R().
					SetHeader("X-API-KEY", cfg.APIKEY).
					SetContentType("application/json").
					SetResult(&response).
					SetError(&resErr).
					Get(APIShortenURL + "?" + query.Encode())
				if err != nil {
					return fmt.Errorf("failed to list short URLs: \n  %w", err)
				}

				if res.StatusCode() != http.StatusOK {
					return fmt.Errorf("failed to list short URLs: \n  status code: %d \n      errcode: %d \n      errinfo: %s",
						res.StatusCode(),
						resErr.ErrCode,
						resErr.ErrInfo)
				}
			}

			// 显示结果
			if len(response.Data) == 0 {
				fmt.Println("No short URLs found")
				return nil
			}

			for _, item := range response.Data {
				fmt.Printf("  Short Code: %s\n", item.Code)
				fmt.Printf("   Short URL: %s\n", item.ShortURL)
				fmt.Printf("Original URL: %s\n", item.OriginalURL)
				if item.Describe != "" {
					fmt.Printf(" Description: %s\n", item.Describe)
				}
				fmt.Println("--------------------------------")
			}

			fmt.Printf("  Total Items: %d\n", response.Meta.TotalItems)
			if !isAll {
				fmt.Printf("  Total Pages: %d\n", response.Meta.TotalPages)
				fmt.Printf("    Page Size: %d\n", response.Meta.PageSize)
				fmt.Printf(" Current Page: %d\n", response.Meta.Page)
				fmt.Printf("Current Count: %d\n", response.Meta.CurrentCount)
			}

			return nil
		},
	}

	cmd.Flags().BoolP("all", "a", false, "List all short links")

	cmd.Flags().Int64P("page", "p", 1, "Page number")
	cmd.Flags().Int64P("psize", "z", 10, "Page size")
	cmd.Flags().StringP("sort", "s", "created_at", "Sort by field")
	cmd.Flags().StringP("order", "o", "asc", "Sort order")

	cmd.Flags().StringP("code", "c", "", "Short code")
	cmd.Flags().StringP("original_url", "r", "", "Original URL")

	return cmd
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(0)
	}
}
