// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sagemaker_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/YakDriver/regexache"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfsagemaker "github.com/hashicorp/terraform-provider-aws/internal/service/sagemaker"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/names"
)

func TestAccSageMakerEndpointConfiguration_basic(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_sagemaker_endpoint_configuration.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.SageMakerServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckEndpointConfigurationDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccEndpointConfigurationConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckEndpointConfigurationExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, names.AttrName, rName),
					acctest.CheckResourceAttrRegionalARNFormat(ctx, resourceName, names.AttrARN, "sagemaker", "endpoint-config/{name}"),
					resource.TestCheckResourceAttr(resourceName, "production_variants.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "production_variants.0.variant_name", "variant-1"),
					resource.TestCheckResourceAttr(resourceName, "production_variants.0.model_name", rName),
					resource.TestCheckResourceAttr(resourceName, "production_variants.0.initial_instance_count", "2"),
					resource.TestCheckResourceAttr(resourceName, "production_variants.0.instance_type", "ml.t2.medium"),
					resource.TestCheckResourceAttr(resourceName, "production_variants.0.initial_variant_weight", "1"),
					resource.TestCheckResourceAttr(resourceName, "production_variants.0.serverless_config.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "production_variants.0.core_dump_config.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "production_variants.0.enable_ssm_access", acctest.CtFalse),
					resource.TestCheckResourceAttr(resourceName, "data_capture_config.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "async_inference_config.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "shadow_production_variants.#", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSageMakerEndpointConfiguration_nameGenerated(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_sagemaker_endpoint_configuration.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.SageMakerServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckEndpointConfigurationDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccEndpointConfigurationConfig_nameGenerated(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckEndpointConfigurationExists(ctx, resourceName),
					acctest.CheckResourceAttrNameGenerated(resourceName, names.AttrName),
					resource.TestCheckResourceAttr(resourceName, names.AttrNamePrefix, "terraform-"),
					acctest.CheckResourceAttrRegionalARNFormat(ctx, resourceName, names.AttrARN, "sagemaker", "endpoint-config/{name}"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSageMakerEndpointConfiguration_namePrefix(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_sagemaker_endpoint_configuration.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.SageMakerServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckEndpointConfigurationDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccEndpointConfigurationConfig_namePrefix(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckEndpointConfigurationExists(ctx, resourceName),
					acctest.CheckResourceAttrNameFromPrefix(resourceName, names.AttrName, rName),
					resource.TestCheckResourceAttr(resourceName, names.AttrNamePrefix, rName),
					acctest.CheckResourceAttrRegionalARNFormat(ctx, resourceName, names.AttrARN, "sagemaker", "endpoint-config/{name}"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSageMakerEndpointConfiguration_shadowProductionVariants(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_sagemaker_endpoint_configuration.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.SageMakerServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckEndpointConfigurationDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccEndpointConfigurationConfig_shadowProductionVariants(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckEndpointConfigurationExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, names.AttrName, rName),
					resource.TestCheckResourceAttr(resourceName, "production_variants.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "production_variants.0.variant_name", "variant-1"),
					resource.TestCheckResourceAttr(resourceName, "production_variants.0.model_name", rName),
					resource.TestCheckResourceAttr(resourceName, "production_variants.0.initial_instance_count", "2"),
					resource.TestCheckResourceAttr(resourceName, "production_variants.0.instance_type", "ml.t2.medium"),
					resource.TestCheckResourceAttr(resourceName, "production_variants.0.initial_variant_weight", "1"),
					resource.TestCheckResourceAttr(resourceName, "production_variants.0.serverless_config.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "shadow_production_variants.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "shadow_production_variants.0.variant_name", "variant-2"),
					resource.TestCheckResourceAttr(resourceName, "shadow_production_variants.0.model_name", rName),
					resource.TestCheckResourceAttr(resourceName, "shadow_production_variants.0.initial_instance_count", "2"),
					resource.TestCheckResourceAttr(resourceName, "shadow_production_variants.0.instance_type", "ml.t2.medium"),
					resource.TestCheckResourceAttr(resourceName, "shadow_production_variants.0.initial_variant_weight", "1"),
					resource.TestCheckResourceAttr(resourceName, "shadow_production_variants.0.serverless_config.#", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSageMakerEndpointConfiguration_ProductionVariants_routing(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_sagemaker_endpoint_configuration.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.SageMakerServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckEndpointConfigurationDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccEndpointConfigurationConfig_routing(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckEndpointConfigurationExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "production_variants.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "production_variants.0.routing_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "production_variants.0.routing_config.0.routing_strategy", "RANDOM"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSageMakerEndpointConfiguration_ProductionVariants_serverless(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_sagemaker_endpoint_configuration.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.SageMakerServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckEndpointConfigurationDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccEndpointConfigurationConfig_serverless(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckEndpointConfigurationExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "production_variants.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "production_variants.0.serverless_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "production_variants.0.serverless_config.0.max_concurrency", "1"),
					resource.TestCheckResourceAttr(resourceName, "production_variants.0.serverless_config.0.memory_size_in_mb", "1024"),
					resource.TestCheckResourceAttr(resourceName, "production_variants.0.serverless_config.0.provisioned_concurrency", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSageMakerEndpointConfiguration_ProductionVariants_ami(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_sagemaker_endpoint_configuration.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.SageMakerServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckEndpointConfigurationDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccEndpointConfigurationConfig_ami(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckEndpointConfigurationExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "production_variants.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "production_variants.0.inference_ami_version", "al2-ami-sagemaker-inference-gpu-2"), //lintignore:AWSAT002
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSageMakerEndpointConfiguration_ProductionVariants_serverlessProvisionedConcurrency(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_sagemaker_endpoint_configuration.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.SageMakerServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckEndpointConfigurationDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccEndpointConfigurationConfig_serverlessProvisionedConcurrency(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckEndpointConfigurationExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "production_variants.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "production_variants.0.serverless_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "production_variants.0.serverless_config.0.max_concurrency", "200"),
					resource.TestCheckResourceAttr(resourceName, "production_variants.0.serverless_config.0.memory_size_in_mb", "5120"),
					resource.TestCheckResourceAttr(resourceName, "production_variants.0.serverless_config.0.provisioned_concurrency", "100"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSageMakerEndpointConfiguration_ProductionVariants_initialVariantWeight(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_sagemaker_endpoint_configuration.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.SageMakerServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckEndpointConfigurationDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccEndpointConfigurationConfig_productionVariantsInitialVariantWeight(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckEndpointConfigurationExists(ctx, resourceName),
					resource.TestCheckResourceAttr(
						resourceName, "production_variants.1.initial_variant_weight", "0.5"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSageMakerEndpointConfiguration_ProductionVariants_acceleratorType(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_sagemaker_endpoint_configuration.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.SageMakerServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckEndpointConfigurationDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccEndpointConfigurationConfig_productionVariantAcceleratorType(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckEndpointConfigurationExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "production_variants.0.accelerator_type", "ml.eia1.medium"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSageMakerEndpointConfiguration_ProductionVariants_variantNameGenerated(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_sagemaker_endpoint_configuration.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.SageMakerServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckEndpointConfigurationDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccEndpointConfigurationConfig_productionVariantVariantNameGenerated(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckEndpointConfigurationExists(ctx, resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "production_variants.0.variant_name"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSageMakerEndpointConfiguration_kmsKeyID(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_sagemaker_endpoint_configuration.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.SageMakerServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckEndpointConfigurationDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccEndpointConfigurationConfig_kmsKeyID(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckEndpointConfigurationExists(ctx, resourceName),
					resource.TestCheckResourceAttrPair(resourceName, names.AttrKMSKeyARN, "aws_kms_key.test", names.AttrARN),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSageMakerEndpointConfiguration_tags(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_sagemaker_endpoint_configuration.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.SageMakerServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckEndpointConfigurationDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccEndpointConfigurationConfig_tags1(rName, acctest.CtKey1, acctest.CtValue1),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckEndpointConfigurationExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, acctest.CtTagsPercent, "1"),
					resource.TestCheckResourceAttr(resourceName, acctest.CtTagsKey1, acctest.CtValue1),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccEndpointConfigurationConfig_tags2(rName, acctest.CtKey1, acctest.CtValue1Updated, acctest.CtKey2, acctest.CtValue2),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckEndpointConfigurationExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, acctest.CtTagsPercent, "2"),
					resource.TestCheckResourceAttr(resourceName, acctest.CtTagsKey1, acctest.CtValue1Updated),
					resource.TestCheckResourceAttr(resourceName, acctest.CtTagsKey2, acctest.CtValue2),
				),
			},
			{
				Config: testAccEndpointConfigurationConfig_tags1(rName, acctest.CtKey2, acctest.CtValue2),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckEndpointConfigurationExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, acctest.CtTagsPercent, "1"),
					resource.TestCheckResourceAttr(resourceName, acctest.CtTagsKey2, acctest.CtValue2),
				),
			},
		},
	})
}

func TestAccSageMakerEndpointConfiguration_dataCapture(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_sagemaker_endpoint_configuration.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.SageMakerServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckEndpointConfigurationDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccEndpointConfigurationConfig_dataCapture(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckEndpointConfigurationExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "data_capture_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "data_capture_config.0.enable_capture", acctest.CtTrue),
					resource.TestCheckResourceAttr(resourceName, "data_capture_config.0.initial_sampling_percentage", "50"),
					resource.TestCheckResourceAttr(resourceName, "data_capture_config.0.destination_s3_uri", fmt.Sprintf("s3://%s/", rName)),
					resource.TestCheckResourceAttr(resourceName, "data_capture_config.0.capture_options.0.capture_mode", "Input"),
					resource.TestCheckResourceAttr(resourceName, "data_capture_config.0.capture_options.1.capture_mode", "Output"),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("data_capture_config").AtSliceIndex(0).AtMapKey("capture_content_type_header"), knownvalue.ListExact([]knownvalue.Check{
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"csv_content_types": knownvalue.Null(),
							"json_content_types": knownvalue.SetExact([]knownvalue.Check{
								knownvalue.StringExact("application/json"),
							}),
						}),
					})),
				},
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSageMakerEndpointConfiguration_dataCapture_inputAndOutput(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_sagemaker_endpoint_configuration.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.SageMakerServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckEndpointConfigurationDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccEndpointConfigurationConfig_dataCapture_inputAndOutput(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckEndpointConfigurationExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "data_capture_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "data_capture_config.0.enable_capture", acctest.CtTrue),
					resource.TestCheckResourceAttr(resourceName, "data_capture_config.0.initial_sampling_percentage", "50"),
					resource.TestCheckResourceAttr(resourceName, "data_capture_config.0.destination_s3_uri", fmt.Sprintf("s3://%s/", rName)),
					resource.TestCheckResourceAttr(resourceName, "data_capture_config.0.capture_options.0.capture_mode", "InputAndOutput"),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("data_capture_config").AtSliceIndex(0).AtMapKey("capture_content_type_header"), knownvalue.ListExact([]knownvalue.Check{
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"csv_content_types": knownvalue.Null(),
							"json_content_types": knownvalue.SetExact([]knownvalue.Check{
								knownvalue.StringExact("application/json"),
							}),
						}),
					})),
				},
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSageMakerEndpointConfiguration_dataCapture_NoHeaders(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_sagemaker_endpoint_configuration.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.SageMakerServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckEndpointConfigurationDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccEndpointConfigurationConfig_dataCapture_noHeaders(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckEndpointConfigurationExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "data_capture_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "data_capture_config.0.enable_capture", acctest.CtTrue),
					resource.TestCheckResourceAttr(resourceName, "data_capture_config.0.initial_sampling_percentage", "50"),
					resource.TestCheckResourceAttr(resourceName, "data_capture_config.0.destination_s3_uri", fmt.Sprintf("s3://%s/", rName)),
					resource.TestCheckResourceAttr(resourceName, "data_capture_config.0.capture_options.0.capture_mode", "Input"),
					resource.TestCheckResourceAttr(resourceName, "data_capture_config.0.capture_options.1.capture_mode", "Output"),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("data_capture_config").AtSliceIndex(0).AtMapKey("capture_content_type_header"), knownvalue.ListExact([]knownvalue.Check{})),
				},
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSageMakerEndpointConfiguration_dataCapture_EmptyHeaders(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.SageMakerServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckEndpointConfigurationDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config:      testAccEndpointConfigurationConfig_dataCapture_emptyHeaders(rName),
				ExpectError: regexache.MustCompile(`At least one attribute out of \[csv_content_types, json_content_types\] must be specified`),
			},
		},
	})
}

func TestAccSageMakerEndpointConfiguration_dataCapture_BothHeaders(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_sagemaker_endpoint_configuration.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.SageMakerServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckEndpointConfigurationDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccEndpointConfigurationConfig_dataCapture_bothHeaders(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckEndpointConfigurationExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "data_capture_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "data_capture_config.0.enable_capture", acctest.CtTrue),
					resource.TestCheckResourceAttr(resourceName, "data_capture_config.0.initial_sampling_percentage", "50"),
					resource.TestCheckResourceAttr(resourceName, "data_capture_config.0.destination_s3_uri", fmt.Sprintf("s3://%s/", rName)),
					resource.TestCheckResourceAttr(resourceName, "data_capture_config.0.capture_options.0.capture_mode", "Input"),
					resource.TestCheckResourceAttr(resourceName, "data_capture_config.0.capture_options.1.capture_mode", "Output"),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("data_capture_config").AtSliceIndex(0).AtMapKey("capture_content_type_header"), knownvalue.ListExact([]knownvalue.Check{
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"csv_content_types": knownvalue.SetExact([]knownvalue.Check{
								knownvalue.StringExact("text/csv"),
							}),
							"json_content_types": knownvalue.SetExact([]knownvalue.Check{
								knownvalue.StringExact("application/json"),
							}),
						}),
					})),
				},
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSageMakerEndpointConfiguration_disappears(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_sagemaker_endpoint_configuration.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.SageMakerServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckEndpointConfigurationDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccEndpointConfigurationConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckEndpointConfigurationExists(ctx, resourceName),
					acctest.CheckResourceDisappears(ctx, acctest.Provider, tfsagemaker.ResourceEndpointConfiguration(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccSageMakerEndpointConfiguration_async(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_sagemaker_endpoint_configuration.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.SageMakerServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckEndpointConfigurationDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccEndpointConfigurationConfig_async(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckEndpointConfigurationExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, names.AttrName, rName),
					resource.TestCheckResourceAttr(resourceName, "async_inference_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "async_inference_config.0.client_config.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "async_inference_config.0.output_config.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "async_inference_config.0.output_config.0.s3_output_path"),
					resource.TestCheckResourceAttr(resourceName, "async_inference_config.0.output_config.0.notification_config.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "async_inference_config.0.output_config.0.kms_key_id", ""),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSageMakerEndpointConfiguration_async_includeInference(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_sagemaker_endpoint_configuration.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.SageMakerServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckEndpointConfigurationDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccEndpointConfigurationConfig_asyncNotifInferenceIn(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckEndpointConfigurationExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, names.AttrName, rName),
					resource.TestCheckResourceAttr(resourceName, "async_inference_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "async_inference_config.0.output_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "async_inference_config.0.output_config.0.notification_config.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "async_inference_config.0.output_config.0.notification_config.0.error_topic", "aws_sns_topic.test", names.AttrARN),
					resource.TestCheckResourceAttrPair(resourceName, "async_inference_config.0.output_config.0.notification_config.0.success_topic", "aws_sns_topic.test", names.AttrARN),
					resource.TestCheckResourceAttr(resourceName, "async_inference_config.0.output_config.0.notification_config.0.include_inference_response_in.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "async_inference_config.0.output_config.0.notification_config.0.include_inference_response_in.*", "SUCCESS_NOTIFICATION_TOPIC"),
					resource.TestCheckTypeSetElemAttr(resourceName, "async_inference_config.0.output_config.0.notification_config.0.include_inference_response_in.*", "ERROR_NOTIFICATION_TOPIC"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSageMakerEndpointConfiguration_async_kms(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_sagemaker_endpoint_configuration.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.SageMakerServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckEndpointConfigurationDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccEndpointConfigurationConfig_asyncKMS(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckEndpointConfigurationExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, names.AttrName, rName),
					resource.TestCheckResourceAttr(resourceName, "async_inference_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "async_inference_config.0.client_config.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "async_inference_config.0.output_config.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "async_inference_config.0.output_config.0.s3_output_path"),
					resource.TestCheckResourceAttr(resourceName, "async_inference_config.0.output_config.0.notification_config.#", "0"),
					resource.TestCheckResourceAttrPair(resourceName, "async_inference_config.0.output_config.0.kms_key_id", "aws_kms_key.test", names.AttrARN),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSageMakerEndpointConfiguration_Async_notif(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_sagemaker_endpoint_configuration.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.SageMakerServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckEndpointConfigurationDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccEndpointConfigurationConfig_asyncNotif(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckEndpointConfigurationExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, names.AttrName, rName),
					resource.TestCheckResourceAttr(resourceName, "async_inference_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "async_inference_config.0.client_config.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "async_inference_config.0.output_config.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "async_inference_config.0.output_config.0.s3_output_path"),
					resource.TestCheckResourceAttr(resourceName, "async_inference_config.0.output_config.0.notification_config.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "async_inference_config.0.output_config.0.notification_config.0.error_topic", "aws_sns_topic.test", names.AttrARN),
					resource.TestCheckResourceAttrPair(resourceName, "async_inference_config.0.output_config.0.notification_config.0.success_topic", "aws_sns_topic.test", names.AttrARN),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSageMakerEndpointConfiguration_Async_client(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_sagemaker_endpoint_configuration.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.SageMakerServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckEndpointConfigurationDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccEndpointConfigurationConfig_asyncClient(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckEndpointConfigurationExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, names.AttrName, rName),
					resource.TestCheckResourceAttr(resourceName, "async_inference_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "async_inference_config.0.client_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "async_inference_config.0.client_config.0.max_concurrent_invocations_per_instance", "1"),
					resource.TestCheckResourceAttr(resourceName, "async_inference_config.0.output_config.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "async_inference_config.0.output_config.0.s3_output_path"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSageMakerEndpointConfiguration_Async_client_failurePath(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_sagemaker_endpoint_configuration.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.SageMakerServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckEndpointConfigurationDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccEndpointConfigurationConfig_asyncClientFailure(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckEndpointConfigurationExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, names.AttrName, rName),
					resource.TestCheckResourceAttr(resourceName, "async_inference_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "async_inference_config.0.client_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "async_inference_config.0.client_config.0.max_concurrent_invocations_per_instance", "1"),
					resource.TestCheckResourceAttr(resourceName, "async_inference_config.0.output_config.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "async_inference_config.0.output_config.0.s3_output_path"),
					resource.TestCheckResourceAttrSet(resourceName, "async_inference_config.0.output_config.0.s3_failure_path"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSageMakerEndpointConfiguration_upgradeToEnableSSMAccess(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_sagemaker_endpoint_configuration.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:   acctest.ErrorCheck(t, names.SageMakerServiceID),
		CheckDestroy: testAccCheckEndpointConfigurationDestroy(ctx),
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"aws": {
						Source:            "hashicorp/aws",
						VersionConstraint: "4.60.0",
					},
				},
				Config: testAccEndpointConfigurationConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckEndpointConfigurationExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, names.AttrName, rName),
					resource.TestCheckResourceAttr(resourceName, "production_variants.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "production_variants.0.variant_name", "variant-1"),
					resource.TestCheckResourceAttr(resourceName, "production_variants.0.model_name", rName),
					resource.TestCheckResourceAttr(resourceName, "production_variants.0.initial_instance_count", "2"),
					resource.TestCheckResourceAttr(resourceName, "production_variants.0.instance_type", "ml.t2.medium"),
					resource.TestCheckResourceAttr(resourceName, "production_variants.0.initial_variant_weight", "1"),
					resource.TestCheckResourceAttr(resourceName, "production_variants.0.serverless_config.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "production_variants.0.core_dump_config.#", "0"),
					resource.TestCheckNoResourceAttr(resourceName, "production_variants.0.enable_ssm_access"),
					resource.TestCheckResourceAttr(resourceName, "data_capture_config.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "async_inference_config.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "shadow_production_variants.#", "0"),
				),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionCreate),
					},
				},
			},
			{
				ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
				Config:                   testAccEndpointConfigurationConfig_basic(rName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionNoop),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionNoop),
					},
				},
			},
		},
	})
}

func TestAccSageMakerEndpointConfiguration_productionVariantsManagedInstanceScaling(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_sagemaker_endpoint_configuration.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.SageMakerServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckEndpointConfigurationDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccEndpointConfigurationConfig_productionVariantsManagedInstanceScaling(rName, 1),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckEndpointConfigurationExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, names.AttrName, rName),
					resource.TestCheckResourceAttr(resourceName, "production_variants.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "production_variants.0.variant_name", "variant-1"),
					resource.TestCheckResourceAttr(resourceName, "production_variants.0.model_name", rName),
					resource.TestCheckResourceAttr(resourceName, "production_variants.0.initial_instance_count", "1"),
					resource.TestCheckResourceAttr(resourceName, "production_variants.0.instance_type", "ml.g5.4xlarge"),
					resource.TestCheckResourceAttr(resourceName, "production_variants.0.managed_instance_scaling.0.status", "ENABLED"),
					resource.TestCheckResourceAttr(resourceName, "production_variants.0.managed_instance_scaling.0.min_instance_count", "1"),
					resource.TestCheckResourceAttr(resourceName, "production_variants.0.managed_instance_scaling.0.max_instance_count", "2"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSageMakerEndpointConfiguration_productionVariantsManagedInstanceScalingZero(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_sagemaker_endpoint_configuration.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.SageMakerServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckEndpointConfigurationDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccEndpointConfigurationConfig_productionVariantsManagedInstanceScaling(rName, 0),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckEndpointConfigurationExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, names.AttrName, rName),
					resource.TestCheckResourceAttr(resourceName, "production_variants.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "production_variants.0.variant_name", "variant-1"),
					resource.TestCheckResourceAttr(resourceName, "production_variants.0.model_name", rName),
					resource.TestCheckResourceAttr(resourceName, "production_variants.0.initial_instance_count", "1"),
					resource.TestCheckResourceAttr(resourceName, "production_variants.0.instance_type", "ml.g5.4xlarge"),
					resource.TestCheckResourceAttr(resourceName, "production_variants.0.managed_instance_scaling.0.status", "ENABLED"),
					resource.TestCheckResourceAttr(resourceName, "production_variants.0.managed_instance_scaling.0.min_instance_count", "0"),
					resource.TestCheckResourceAttr(resourceName, "production_variants.0.managed_instance_scaling.0.max_instance_count", "2"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckEndpointConfigurationDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).SageMakerClient(ctx)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_sagemaker_endpoint_configuration" {
				continue
			}

			_, err := tfsagemaker.FindEndpointConfigByName(ctx, conn, rs.Primary.ID)

			if tfresource.NotFound(err) {
				continue
			}

			if err != nil {
				return err
			}

			return fmt.Errorf("SageMaker AI Endpoint Configuration %s still exists", rs.Primary.ID)
		}

		return nil
	}
}

func testAccCheckEndpointConfigurationExists(ctx context.Context, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("SageMaker AI endpoint config not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no SageMaker AI endpoint config ID is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).SageMakerClient(ctx)
		_, err := tfsagemaker.FindEndpointConfigByName(ctx, conn, rs.Primary.ID)

		return err
	}
}

func testAccEndpointConfigurationConfig_base(rName string) string {
	return fmt.Sprintf(`
data "aws_sagemaker_prebuilt_ecr_image" "test" {
  repository_name = "kmeans"
}

resource "aws_sagemaker_model" "test" {
  name               = %[1]q
  execution_role_arn = aws_iam_role.test.arn

  primary_container {
    image = data.aws_sagemaker_prebuilt_ecr_image.test.registry_path
  }
}

resource "aws_iam_role" "test" {
  name               = %[1]q
  path               = "/"
  assume_role_policy = data.aws_iam_policy_document.assume_role.json
}

data "aws_iam_policy_document" "assume_role" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["sagemaker.amazonaws.com"]
    }
  }
}
`, rName)
}

func testAccEndpointConfigurationConfig_basic(rName string) string {
	return acctest.ConfigCompose(testAccEndpointConfigurationConfig_base(rName), fmt.Sprintf(`
resource "aws_sagemaker_endpoint_configuration" "test" {
  name = %[1]q

  production_variants {
    variant_name           = "variant-1"
    model_name             = aws_sagemaker_model.test.name
    initial_instance_count = 2
    instance_type          = "ml.t2.medium"
    initial_variant_weight = 1
  }
}
`, rName))
}

func testAccEndpointConfigurationConfig_nameGenerated(rName string) string {
	return acctest.ConfigCompose(testAccEndpointConfigurationConfig_base(rName), `
resource "aws_sagemaker_endpoint_configuration" "test" {
  production_variants {
    variant_name           = "variant-1"
    model_name             = aws_sagemaker_model.test.name
    initial_instance_count = 2
    instance_type          = "ml.t2.medium"
    initial_variant_weight = 1
  }
}
`)
}

func testAccEndpointConfigurationConfig_namePrefix(rName string) string {
	return acctest.ConfigCompose(testAccEndpointConfigurationConfig_base(rName), fmt.Sprintf(`
resource "aws_sagemaker_endpoint_configuration" "test" {
  name_prefix = %[1]q

  production_variants {
    variant_name           = "variant-1"
    model_name             = aws_sagemaker_model.test.name
    initial_instance_count = 2
    instance_type          = "ml.t2.medium"
    initial_variant_weight = 1
  }
}
`, rName))
}

func testAccEndpointConfigurationConfig_shadowProductionVariants(rName string) string {
	return acctest.ConfigCompose(testAccEndpointConfigurationConfig_base(rName), fmt.Sprintf(`
resource "aws_sagemaker_endpoint_configuration" "test" {
  name = %[1]q

  production_variants {
    variant_name           = "variant-1"
    model_name             = aws_sagemaker_model.test.name
    initial_instance_count = 2
    instance_type          = "ml.t2.medium"
    initial_variant_weight = 1
  }

  shadow_production_variants {
    variant_name           = "variant-2"
    model_name             = aws_sagemaker_model.test.name
    initial_instance_count = 2
    instance_type          = "ml.t2.medium"
    initial_variant_weight = 1
  }
}
`, rName))
}

func testAccEndpointConfigurationConfig_productionVariantsInitialVariantWeight(rName string) string {
	return acctest.ConfigCompose(testAccEndpointConfigurationConfig_base(rName), fmt.Sprintf(`
resource "aws_sagemaker_endpoint_configuration" "test" {
  name = %[1]q

  production_variants {
    variant_name           = "variant-1"
    model_name             = aws_sagemaker_model.test.name
    initial_instance_count = 1
    instance_type          = "ml.t2.medium"
  }

  production_variants {
    variant_name           = "variant-2"
    model_name             = aws_sagemaker_model.test.name
    initial_instance_count = 1
    instance_type          = "ml.t2.medium"
    initial_variant_weight = 0.5
  }
}
`, rName))
}

func testAccEndpointConfigurationConfig_productionVariantAcceleratorType(rName string) string {
	return acctest.ConfigCompose(testAccEndpointConfigurationConfig_base(rName), fmt.Sprintf(`
resource "aws_sagemaker_endpoint_configuration" "test" {
  name = %[1]q

  production_variants {
    variant_name           = "variant-1"
    model_name             = aws_sagemaker_model.test.name
    initial_instance_count = 2
    instance_type          = "ml.t2.medium"
    accelerator_type       = "ml.eia1.medium"
    initial_variant_weight = 1
  }
}
`, rName))
}

func testAccEndpointConfigurationConfig_productionVariantVariantNameGenerated(rName string) string {
	return acctest.ConfigCompose(testAccEndpointConfigurationConfig_base(rName), fmt.Sprintf(`
resource "aws_sagemaker_endpoint_configuration" "test" {
  name = %[1]q

  production_variants {
    model_name             = aws_sagemaker_model.test.name
    initial_instance_count = 2
    instance_type          = "ml.t2.medium"
    initial_variant_weight = 1
  }
}
`, rName))
}

func testAccEndpointConfigurationConfig_kmsKeyID(rName string) string {
	return acctest.ConfigCompose(testAccEndpointConfigurationConfig_base(rName), fmt.Sprintf(`
resource "aws_sagemaker_endpoint_configuration" "test" {
  name        = %[1]q
  kms_key_arn = aws_kms_key.test.arn

  production_variants {
    variant_name           = "variant-1"
    model_name             = aws_sagemaker_model.test.name
    initial_instance_count = 1
    instance_type          = "ml.t2.medium"
    initial_variant_weight = 1
  }
}

resource "aws_kms_key" "test" {
  description             = %[1]q
  deletion_window_in_days = 10
  enable_key_rotation     = true
}
`, rName))
}

func testAccEndpointConfigurationConfig_tags1(rName, tagKey1, tagValue1 string) string {
	return acctest.ConfigCompose(testAccEndpointConfigurationConfig_base(rName), fmt.Sprintf(`
resource "aws_sagemaker_endpoint_configuration" "test" {
  name = %[1]q

  production_variants {
    variant_name           = "variant-1"
    model_name             = aws_sagemaker_model.test.name
    initial_instance_count = 1
    instance_type          = "ml.t2.medium"
    initial_variant_weight = 1
  }

  tags = {
    %[2]q = %[3]q
  }
}
`, rName, tagKey1, tagValue1))
}

func testAccEndpointConfigurationConfig_tags2(rName, tagKey1, tagValue1, tagKey2, tagValue2 string) string {
	return acctest.ConfigCompose(testAccEndpointConfigurationConfig_base(rName), fmt.Sprintf(`
resource "aws_sagemaker_endpoint_configuration" "test" {
  name = %[1]q

  production_variants {
    variant_name           = "variant-1"
    model_name             = aws_sagemaker_model.test.name
    initial_instance_count = 1
    instance_type          = "ml.t2.medium"
    initial_variant_weight = 1
  }

  tags = {
    %[2]q = %[3]q
    %[4]q = %[5]q
  }
}
`, rName, tagKey1, tagValue1, tagKey2, tagValue2))
}

func testAccEndpointConfigurationConfig_dataCapture(rName string) string {
	return acctest.ConfigCompose(testAccEndpointConfigurationConfig_base(rName), fmt.Sprintf(`
resource "aws_s3_bucket" "test" {
  bucket        = %[1]q
  force_destroy = true
}

resource "aws_sagemaker_endpoint_configuration" "test" {
  name = %[1]q

  production_variants {
    variant_name           = "variant-1"
    model_name             = aws_sagemaker_model.test.name
    initial_instance_count = 2
    instance_type          = "ml.t2.medium"
    initial_variant_weight = 1
  }

  data_capture_config {
    enable_capture              = true
    initial_sampling_percentage = 50
    destination_s3_uri          = "s3://${aws_s3_bucket.test.bucket}/"

    capture_options {
      capture_mode = "Input"
    }

    capture_options {
      capture_mode = "Output"
    }

    capture_content_type_header {
      json_content_types = ["application/json"]
    }
  }
}
`, rName))
}

func testAccEndpointConfigurationConfig_dataCapture_inputAndOutput(rName string) string {
	return acctest.ConfigCompose(testAccEndpointConfigurationConfig_base(rName), fmt.Sprintf(`
resource "aws_s3_bucket" "test" {
  bucket        = %[1]q
  force_destroy = true
}

resource "aws_sagemaker_endpoint_configuration" "test" {
  name = %[1]q

  production_variants {
    variant_name           = "variant-1"
    model_name             = aws_sagemaker_model.test.name
    initial_instance_count = 2
    instance_type          = "ml.t2.medium"
    initial_variant_weight = 1
  }

  data_capture_config {
    enable_capture              = true
    initial_sampling_percentage = 50
    destination_s3_uri          = "s3://${aws_s3_bucket.test.bucket}/"

    capture_options {
      capture_mode = "InputAndOutput"
    }

    capture_content_type_header {
      json_content_types = ["application/json"]
    }
  }
}
`, rName))
}

func testAccEndpointConfigurationConfig_dataCapture_noHeaders(rName string) string {
	return acctest.ConfigCompose(testAccEndpointConfigurationConfig_base(rName), fmt.Sprintf(`
resource "aws_s3_bucket" "test" {
  bucket        = %[1]q
  force_destroy = true
}

resource "aws_sagemaker_endpoint_configuration" "test" {
  name = %[1]q

  production_variants {
    variant_name           = "variant-1"
    model_name             = aws_sagemaker_model.test.name
    initial_instance_count = 2
    instance_type          = "ml.t2.medium"
    initial_variant_weight = 1
  }

  data_capture_config {
    enable_capture              = true
    initial_sampling_percentage = 50
    destination_s3_uri          = "s3://${aws_s3_bucket.test.bucket}/"

    capture_options {
      capture_mode = "Input"
    }

    capture_options {
      capture_mode = "Output"
    }
  }
}
`, rName))
}

func testAccEndpointConfigurationConfig_dataCapture_emptyHeaders(rName string) string {
	return acctest.ConfigCompose(testAccEndpointConfigurationConfig_base(rName), fmt.Sprintf(`
resource "aws_s3_bucket" "test" {
  bucket        = %[1]q
  force_destroy = true
}

resource "aws_sagemaker_endpoint_configuration" "test" {
  name = %[1]q

  production_variants {
    variant_name           = "variant-1"
    model_name             = aws_sagemaker_model.test.name
    initial_instance_count = 2
    instance_type          = "ml.t2.medium"
    initial_variant_weight = 1
  }

  data_capture_config {
    enable_capture              = true
    initial_sampling_percentage = 50
    destination_s3_uri          = "s3://${aws_s3_bucket.test.bucket}/"

    capture_options {
      capture_mode = "Input"
    }

    capture_options {
      capture_mode = "Output"
    }

    capture_content_type_header {
    }
  }
}
`, rName))
}

func testAccEndpointConfigurationConfig_dataCapture_bothHeaders(rName string) string {
	return acctest.ConfigCompose(testAccEndpointConfigurationConfig_base(rName), fmt.Sprintf(`
resource "aws_s3_bucket" "test" {
  bucket        = %[1]q
  force_destroy = true
}

resource "aws_sagemaker_endpoint_configuration" "test" {
  name = %[1]q

  production_variants {
    variant_name           = "variant-1"
    model_name             = aws_sagemaker_model.test.name
    initial_instance_count = 2
    instance_type          = "ml.t2.medium"
    initial_variant_weight = 1
  }

  data_capture_config {
    enable_capture              = true
    initial_sampling_percentage = 50
    destination_s3_uri          = "s3://${aws_s3_bucket.test.bucket}/"

    capture_options {
      capture_mode = "Input"
    }

    capture_options {
      capture_mode = "Output"
    }

    capture_content_type_header {
      csv_content_types  = ["text/csv"]
      json_content_types = ["application/json"]
    }
  }
}
`, rName))
}

func testAccEndpointConfigurationConfig_asyncKMS(rName string) string {
	return acctest.ConfigCompose(testAccEndpointConfigurationConfig_base(rName), fmt.Sprintf(`
resource "aws_s3_bucket" "test" {
  bucket        = %[1]q
  force_destroy = true
}

resource "aws_kms_key" "test" {
  description             = %[1]q
  deletion_window_in_days = 7
  enable_key_rotation     = true
}

resource "aws_sagemaker_endpoint_configuration" "test" {
  name = %[1]q

  production_variants {
    variant_name           = "variant-1"
    model_name             = aws_sagemaker_model.test.name
    initial_instance_count = 2
    instance_type          = "ml.t2.medium"
    initial_variant_weight = 1
  }

  async_inference_config {
    output_config {
      s3_output_path = "s3://${aws_s3_bucket.test.bucket}/"
      kms_key_id     = aws_kms_key.test.arn
    }
  }
}
`, rName))
}

func testAccEndpointConfigurationConfig_async(rName string) string {
	return acctest.ConfigCompose(testAccEndpointConfigurationConfig_base(rName), fmt.Sprintf(`
resource "aws_s3_bucket" "test" {
  bucket        = %[1]q
  acl           = "private"
  force_destroy = true
}

resource "aws_sagemaker_endpoint_configuration" "test" {
  name = %[1]q

  production_variants {
    variant_name           = "variant-1"
    model_name             = aws_sagemaker_model.test.name
    initial_instance_count = 2
    instance_type          = "ml.t2.medium"
    initial_variant_weight = 1
  }

  async_inference_config {
    output_config {
      s3_output_path = "s3://${aws_s3_bucket.test.bucket}/"
    }
  }
}
`, rName))
}

func testAccEndpointConfigurationConfig_asyncNotif(rName string) string {
	return acctest.ConfigCompose(testAccEndpointConfigurationConfig_base(rName), fmt.Sprintf(`
resource "aws_s3_bucket" "test" {
  bucket        = %[1]q
  force_destroy = true
}

resource "aws_sns_topic" "test" {
  name = %[1]q
}

resource "aws_kms_key" "test" {
  description             = %[1]q
  deletion_window_in_days = 7
  enable_key_rotation     = true
}

resource "aws_sagemaker_endpoint_configuration" "test" {
  name = %[1]q

  production_variants {
    variant_name           = "variant-1"
    model_name             = aws_sagemaker_model.test.name
    initial_instance_count = 2
    instance_type          = "ml.t2.medium"
    initial_variant_weight = 1
  }

  async_inference_config {
    output_config {
      s3_output_path = "s3://${aws_s3_bucket.test.bucket}/"
      kms_key_id     = aws_kms_key.test.arn

      notification_config {
        error_topic   = aws_sns_topic.test.arn
        success_topic = aws_sns_topic.test.arn
      }
    }
  }
}
`, rName))
}

func testAccEndpointConfigurationConfig_asyncNotifInferenceIn(rName string) string {
	return acctest.ConfigCompose(testAccEndpointConfigurationConfig_base(rName), fmt.Sprintf(`
resource "aws_s3_bucket" "test" {
  bucket        = %[1]q
  force_destroy = true
}

resource "aws_sns_topic" "test" {
  name = %[1]q
}

resource "aws_kms_key" "test" {
  description             = %[1]q
  deletion_window_in_days = 7
  enable_key_rotation     = true
}

resource "aws_sagemaker_endpoint_configuration" "test" {
  name = %[1]q

  production_variants {
    variant_name           = "variant-1"
    model_name             = aws_sagemaker_model.test.name
    initial_instance_count = 2
    instance_type          = "ml.t2.medium"
    initial_variant_weight = 1
  }

  async_inference_config {
    output_config {
      s3_output_path = "s3://${aws_s3_bucket.test.bucket}/"
      kms_key_id     = aws_kms_key.test.arn

      notification_config {
        error_topic                   = aws_sns_topic.test.arn
        include_inference_response_in = ["SUCCESS_NOTIFICATION_TOPIC", "ERROR_NOTIFICATION_TOPIC"]
        success_topic                 = aws_sns_topic.test.arn
      }
    }
  }
}
`, rName))
}

func testAccEndpointConfigurationConfig_asyncClient(rName string) string {
	return acctest.ConfigCompose(testAccEndpointConfigurationConfig_base(rName), fmt.Sprintf(`
resource "aws_s3_bucket" "test" {
  bucket        = %[1]q
  force_destroy = true
}

resource "aws_kms_key" "test" {
  description             = %[1]q
  deletion_window_in_days = 7
  enable_key_rotation     = true
}

resource "aws_sagemaker_endpoint_configuration" "test" {
  name = %[1]q

  production_variants {
    variant_name           = "variant-1"
    model_name             = aws_sagemaker_model.test.name
    initial_instance_count = 2
    instance_type          = "ml.t2.medium"
    initial_variant_weight = 1
  }

  async_inference_config {
    client_config {
      max_concurrent_invocations_per_instance = 1
    }

    output_config {
      s3_output_path = "s3://${aws_s3_bucket.test.bucket}/"
      kms_key_id     = aws_kms_key.test.arn
    }
  }
}
`, rName))
}

func testAccEndpointConfigurationConfig_asyncClientFailure(rName string) string {
	return acctest.ConfigCompose(testAccEndpointConfigurationConfig_base(rName), fmt.Sprintf(`
resource "aws_s3_bucket" "test" {
  bucket        = %[1]q
  force_destroy = true
}

resource "aws_kms_key" "test" {
  description             = %[1]q
  deletion_window_in_days = 7
  enable_key_rotation     = true
}

resource "aws_sagemaker_endpoint_configuration" "test" {
  name = %[1]q

  production_variants {
    variant_name           = "variant-1"
    model_name             = aws_sagemaker_model.test.name
    initial_instance_count = 2
    instance_type          = "ml.t2.medium"
    initial_variant_weight = 1
  }

  async_inference_config {
    client_config {
      max_concurrent_invocations_per_instance = 1
    }

    output_config {
      s3_output_path  = "s3://${aws_s3_bucket.test.bucket}/"
      s3_failure_path = "s3://${aws_s3_bucket.test.bucket}/"
      kms_key_id      = aws_kms_key.test.arn
    }
  }
}
`, rName))
}

func testAccEndpointConfigurationConfig_routing(rName string) string {
	return acctest.ConfigCompose(testAccEndpointConfigurationConfig_base(rName), fmt.Sprintf(`
resource "aws_sagemaker_endpoint_configuration" "test" {
  name = %[1]q

  production_variants {
    variant_name           = "variant-1"
    model_name             = aws_sagemaker_model.test.name
    initial_instance_count = 2
    instance_type          = "ml.t2.medium"

    routing_config {
      routing_strategy = "RANDOM"
    }
  }
}
`, rName))
}

func testAccEndpointConfigurationConfig_serverless(rName string) string {
	return acctest.ConfigCompose(testAccEndpointConfigurationConfig_base(rName), fmt.Sprintf(`
resource "aws_sagemaker_endpoint_configuration" "test" {
  name = %[1]q

  production_variants {
    variant_name = "variant-1"
    model_name   = aws_sagemaker_model.test.name

    serverless_config {
      max_concurrency   = 1
      memory_size_in_mb = 1024
    }
  }
}
`, rName))
}

func testAccEndpointConfigurationConfig_ami(rName string) string {
	//lintignore:AWSAT002
	return acctest.ConfigCompose(testAccEndpointConfigurationConfig_base(rName), fmt.Sprintf(`
resource "aws_sagemaker_endpoint_configuration" "test" {
  name = %[1]q

  production_variants {
    variant_name           = "variant-1"
    model_name             = aws_sagemaker_model.test.name
    inference_ami_version  = "al2-ami-sagemaker-inference-gpu-2"
    instance_type          = "ml.t2.medium"
    initial_instance_count = 2
    initial_variant_weight = 1
  }
}
`, rName))
}

func testAccEndpointConfigurationConfig_serverlessProvisionedConcurrency(rName string) string {
	return acctest.ConfigCompose(testAccEndpointConfigurationConfig_base(rName), fmt.Sprintf(`
resource "aws_sagemaker_endpoint_configuration" "test" {
  name = %[1]q

  production_variants {
    variant_name = "variant-1"
    model_name   = aws_sagemaker_model.test.name

    serverless_config {
      max_concurrency         = 200
      memory_size_in_mb       = 5120
      provisioned_concurrency = 100
    }
  }
}
`, rName))
}

func testAccEndpointConfigurationConfig_productionVariantsManagedInstanceScaling(rName string, min int) string {
	return acctest.ConfigCompose(fmt.Sprintf(`
data "aws_region" "current" {}
data "aws_partition" "current" {}
data "aws_sagemaker_prebuilt_ecr_image" "managed_instance_scaling_test" {
  repository_name = "djl-inference"
  image_tag       = "0.27.0-deepspeed0.12.6-cu121"
}

data "aws_iam_policy_document" "managed_instance_scaling_test_policy" {
  statement {
    effect = "Allow"

    actions = [
      "cloudwatch:PutMetricData",
      "logs:CreateLogStream",
      "logs:PutLogEvents",
      "logs:CreateLogGroup",
      "logs:DescribeLogStreams",
      "ecr:GetAuthorizationToken",
      "ecr:BatchCheckLayerAvailability",
      "ecr:GetDownloadUrlForLayer",
      "ecr:BatchGetImage",
    ]

    resources = [
      "*",
    ]
  }

  statement {
    effect = "Allow"

    actions = [
      "s3:GetObject",
      "s3:ListBucket",
    ]

    resources = [
      aws_s3_bucket.managed_instance_scaling_test.arn,
      "${aws_s3_bucket.managed_instance_scaling_test.arn}/*",
    ]
  }
}

resource "aws_iam_policy" "managed_instance_scaling_test" {
  name        = %[1]q
  description = "Allow SageMaker AI to create model"
  policy      = data.aws_iam_policy_document.managed_instance_scaling_test_policy.json
}

resource "aws_iam_role" "managed_instance_scaling_test" {
  name               = %[1]q
  path               = "/"
  assume_role_policy = data.aws_iam_policy_document.assume_role.json
}

data "aws_iam_policy_document" "assume_role" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["sagemaker.amazonaws.com"]
    }
  }
}

resource "aws_iam_role_policy_attachment" "managed_instance_scaling_test" {
  role       = aws_iam_role.managed_instance_scaling_test.name
  policy_arn = aws_iam_policy.managed_instance_scaling_test.arn
}

resource "aws_s3_bucket" "managed_instance_scaling_test" {
  bucket        = %[1]q
  force_destroy = true
}

resource "aws_s3_object" "managed_instance_scaling_test" {
  bucket  = aws_s3_bucket.managed_instance_scaling_test.bucket
  key     = "model/inference.py"
  content = "some-data"
}

resource "aws_sagemaker_model" "managed_instance_scaling_test" {
  name               = %[1]q
  execution_role_arn = aws_iam_role.managed_instance_scaling_test.arn
  primary_container {
    image = data.aws_sagemaker_prebuilt_ecr_image.managed_instance_scaling_test.registry_path
    model_data_source {
      s3_data_source {
        s3_data_type     = "S3Prefix"
        s3_uri           = "s3://${aws_s3_object.managed_instance_scaling_test.bucket}/model/"
        compression_type = "None"
      }
    }
  }
  depends_on = [
    aws_iam_role_policy_attachment.managed_instance_scaling_test
  ]
}

resource "aws_sagemaker_endpoint_configuration" "test" {
  name = %[1]q

  production_variants {
    variant_name           = "variant-1"
    model_name             = aws_sagemaker_model.managed_instance_scaling_test.name
    initial_instance_count = 1
    instance_type          = "ml.g5.4xlarge"

    managed_instance_scaling {
      status             = "ENABLED"
      min_instance_count = %[2]d
      max_instance_count = 2
    }

    routing_config {
      routing_strategy = "LEAST_OUTSTANDING_REQUESTS"
    }

    model_data_download_timeout_in_seconds            = 60
    container_startup_health_check_timeout_in_seconds = 60
  }
}
`, rName, min))
}
