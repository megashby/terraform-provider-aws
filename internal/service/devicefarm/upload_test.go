// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package devicefarm_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/YakDriver/regexache"
	awstypes "github.com/aws/aws-sdk-go-v2/service/devicefarm/types"
	"github.com/hashicorp/aws-sdk-go-base/v2/endpoints"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	tfknownvalue "github.com/hashicorp/terraform-provider-aws/internal/acctest/knownvalue"
	tfstatecheck "github.com/hashicorp/terraform-provider-aws/internal/acctest/statecheck"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfdevicefarm "github.com/hashicorp/terraform-provider-aws/internal/service/devicefarm"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/names"
)

func TestAccDeviceFarmUpload_basic(t *testing.T) {
	ctx := acctest.Context(t)
	var proj awstypes.Upload
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	rNameUpdated := sdkacctest.RandomWithPrefix("tf-acc-test-updated")
	resourceName := "aws_devicefarm_upload.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionHasService(t, names.DeviceFarmEndpointID)
			// Currently, DeviceFarm is only supported in us-west-2
			// https://docs.aws.amazon.com/general/latest/gr/devicefarm.html
			acctest.PreCheckRegion(t, endpoints.UsWest2RegionID)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.DeviceFarmServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckUploadDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccUploadConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUploadExists(ctx, resourceName, &proj),
					resource.TestCheckResourceAttr(resourceName, names.AttrName, rName),
					acctest.MatchResourceAttrRegionalARN(ctx, resourceName, names.AttrARN, "devicefarm", regexache.MustCompile(`upload:.+`)),
					resource.TestCheckResourceAttr(resourceName, names.AttrType, "APPIUM_JAVA_TESTNG_TEST_SPEC"),
					resource.TestCheckResourceAttr(resourceName, "category", "PRIVATE"),
					resource.TestCheckResourceAttrSet(resourceName, names.AttrURL),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{names.AttrURL},
			},
			{
				Config: testAccUploadConfig_basic(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUploadExists(ctx, resourceName, &proj),
					resource.TestCheckResourceAttr(resourceName, names.AttrName, rNameUpdated),
					acctest.MatchResourceAttrRegionalARN(ctx, resourceName, names.AttrARN, "devicefarm", regexache.MustCompile(`upload:.+`)),
					resource.TestCheckResourceAttr(resourceName, names.AttrType, "APPIUM_JAVA_TESTNG_TEST_SPEC"),
					resource.TestCheckResourceAttr(resourceName, "category", "PRIVATE"),
					resource.TestCheckResourceAttrSet(resourceName, names.AttrURL),
				),
			},
		},
	})
}

func TestAccDeviceFarmUpload_disappears(t *testing.T) {
	ctx := acctest.Context(t)
	var proj awstypes.Upload
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_devicefarm_upload.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionHasService(t, names.DeviceFarmEndpointID)
			// Currently, DeviceFarm is only supported in us-west-2
			// https://docs.aws.amazon.com/general/latest/gr/devicefarm.html
			acctest.PreCheckRegion(t, endpoints.UsWest2RegionID)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.DeviceFarmServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckUploadDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccUploadConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUploadExists(ctx, resourceName, &proj),
					acctest.CheckResourceDisappears(ctx, acctest.Provider, tfdevicefarm.ResourceUpload(), resourceName),
					acctest.CheckResourceDisappears(ctx, acctest.Provider, tfdevicefarm.ResourceUpload(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccDeviceFarmUpload_disappears_project(t *testing.T) {
	ctx := acctest.Context(t)
	var proj awstypes.Upload
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_devicefarm_upload.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionHasService(t, names.DeviceFarmEndpointID)
			// Currently, DeviceFarm is only supported in us-west-2
			// https://docs.aws.amazon.com/general/latest/gr/devicefarm.html
			acctest.PreCheckRegion(t, endpoints.UsWest2RegionID)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.DeviceFarmServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckUploadDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccUploadConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUploadExists(ctx, resourceName, &proj),
					acctest.CheckResourceDisappears(ctx, acctest.Provider, tfdevicefarm.ResourceProject(), "aws_devicefarm_project.test"),
					acctest.CheckResourceDisappears(ctx, acctest.Provider, tfdevicefarm.ResourceUpload(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccDeviceFarmUpload_Identity_ExistingResource(t *testing.T) {
	ctx := acctest.Context(t)
	var proj awstypes.Upload
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_devicefarm_upload.test"

	resource.ParallelTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_12_0),
		},
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionHasService(t, names.DeviceFarmEndpointID)
			// Currently, DeviceFarm is only supported in us-west-2
			// https://docs.aws.amazon.com/general/latest/gr/devicefarm.html
			acctest.PreCheckRegion(t, endpoints.UsWest2RegionID)
		},
		ErrorCheck:   acctest.ErrorCheck(t, names.DeviceFarmServiceID),
		CheckDestroy: testAccCheckUploadDestroy(ctx),
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"aws": {
						Source:            "hashicorp/aws",
						VersionConstraint: "5.100.0",
					},
				},
				Config: testAccUploadConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUploadExists(ctx, resourceName, &proj),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					tfstatecheck.ExpectNoIdentity(resourceName),
				},
			},
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"aws": {
						Source:            "hashicorp/aws",
						VersionConstraint: "6.0.0",
					},
				},
				Config: testAccUploadConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUploadExists(ctx, resourceName, &proj),
				),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionNoop),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionNoop),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectIdentity(resourceName, map[string]knownvalue.Check{
						names.AttrARN: knownvalue.Null(),
					}),
				},
			},
			{
				ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
				Config:                   testAccUploadConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUploadExists(ctx, resourceName, &proj),
				),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionNoop),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionNoop),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectIdentity(resourceName, map[string]knownvalue.Check{
						names.AttrARN: tfknownvalue.RegionalARNRegexp("devicefarm", regexache.MustCompile(`upload:.+`)),
					}),
				},
			},
		},
	})
}

func testAccCheckUploadExists(ctx context.Context, n string, v *awstypes.Upload) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).DeviceFarmClient(ctx)
		resp, err := tfdevicefarm.FindUploadByARN(ctx, conn, rs.Primary.ID)
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("DeviceFarm Upload not found")
		}

		*v = *resp

		return nil
	}
}

func testAccCheckUploadDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).DeviceFarmClient(ctx)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_devicefarm_upload" {
				continue
			}

			// Try to find the resource
			_, err := tfdevicefarm.FindUploadByARN(ctx, conn, rs.Primary.ID)
			if tfresource.NotFound(err) {
				continue
			}

			if err != nil {
				return err
			}

			return fmt.Errorf("DeviceFarm Upload %s still exists", rs.Primary.ID)
		}

		return nil
	}
}

func testAccUploadConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "aws_devicefarm_project" "test" {
  name = %[1]q
}

resource "aws_devicefarm_upload" "test" {
  name        = %[1]q
  project_arn = aws_devicefarm_project.test.arn
  type        = "APPIUM_JAVA_TESTNG_TEST_SPEC"
}
`, rName)
}
