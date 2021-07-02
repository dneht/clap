/*
Copyright 2020 Dasheng.

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

package base

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/utils"
	"time"
)

var App *fiber.App

func WebInit() {
	App = fiber.New()
	App.Use(
		logger.New(logger.Config{
			Format:       "[http] [${status}] [${method}] ${time} - ${path} -- ${ip}, latency=${latency}, bytes_in=${bytesSent}, bytes_out=${bytesReceived}, error=${error}\n",
			TimeFormat:   "2006-01-02 15:04:05",
			TimeZone:     "Local",
			TimeInterval: 5 * time.Second,
		}),
		compress.New(compress.Config{
			Level: compress.LevelBestSpeed,
		}),
		cors.New(cors.Config{
			AllowOrigins: "*",
			AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		}))
	if !IsOffline() {
		App.Use(
			csrf.New(csrf.Config{
				KeyLookup:      "header:X-Csrf-Token",
				CookieName:     "csrf_clap",
				CookieSameSite: "None",
				Expiration:     1 * time.Hour,
				KeyGenerator:   utils.UUID,
			}),
			func(c *fiber.Ctx) error {
				c.Set("X-XSS-Protection", "1; mode=block")
				c.Set("X-Content-Type-Options", "nosniff")
				c.Set("X-Download-Options", "noopen")
				c.Set("Strict-Transport-Security", "max-age=5184000")
				c.Set("X-Frame-Options", "SAMEORIGIN")
				c.Set("X-DNS-Prefetch-Control", "off")
				return c.Next()
			})
	}
}
