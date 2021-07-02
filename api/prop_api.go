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

package api

import (
	"cana.io/clap/pkg/base"
	"cana.io/clap/pkg/refer"
	"github.com/gofiber/fiber/v2"
)

func PropBase(c *fiber.Ctx) error {
	now := base.Now()
	return c.JSON(map[string]interface{}{
		"env":       now.Env,
		"type":      refer.AppTypeMap,
		"namespace": now.Namespace,
		"timezone":  now.Timezone,
		"document":  now.Document,
		"package":   now.Package,
	})
}
