package medialive_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/medialive"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/create"
	tfmedialive "github.com/hashicorp/terraform-provider-aws/internal/service/medialive"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/names"
)

func testAccMultiplex_basic(t *testing.T) {
	ctx := acctest.Context(t)
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	var multiplex medialive.DescribeMultiplexOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_medialive_multiplex.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPartitionHasService(names.MediaLiveEndpointID, t)
			testAccMultiplexesPreCheck(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.MediaLiveEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckMultiplexDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccMultiplexConfig_basic(rName, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMultiplexExists(ctx, resourceName, &multiplex),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrSet(resourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "multiplex_settings.0.transport_stream_bitrate", "1000000"),
					resource.TestCheckResourceAttr(resourceName, "multiplex_settings.0.transport_stream_reserved_bitrate", "1"),
					resource.TestCheckResourceAttr(resourceName, "multiplex_settings.0.transport_stream_id", "1"),
					resource.TestCheckResourceAttr(resourceName, "multiplex_settings.0.maximum_video_buffer_delay_milliseconds", "1000"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"start_multiplex"},
			},
		},
	})
}

func testAccMultiplex_start(t *testing.T) {
	ctx := acctest.Context(t)
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	var multiplex medialive.DescribeMultiplexOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_medialive_multiplex.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPartitionHasService(names.MediaLiveEndpointID, t)
			testAccMultiplexesPreCheck(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.MediaLiveEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckMultiplexDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccMultiplexConfig_basic(rName, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMultiplexExists(ctx, resourceName, &multiplex),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrSet(resourceName, "arn"),
				),
			},
			{
				Config: testAccMultiplexConfig_basic(rName, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMultiplexExists(ctx, resourceName, &multiplex),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrSet(resourceName, "arn"),
				),
			},
		},
	})
}

func testAccMultiplex_update(t *testing.T) {
	ctx := acctest.Context(t)
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	var multiplex medialive.DescribeMultiplexOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_medialive_multiplex.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPartitionHasService(names.MediaLiveEndpointID, t)
			testAccMultiplexesPreCheck(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.MediaLiveEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckMultiplexDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccMultiplexConfig_basic(rName, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMultiplexExists(ctx, resourceName, &multiplex),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrSet(resourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "multiplex_settings.0.transport_stream_bitrate", "1000000"),
					resource.TestCheckResourceAttr(resourceName, "multiplex_settings.0.transport_stream_reserved_bitrate", "1"),
					resource.TestCheckResourceAttr(resourceName, "multiplex_settings.0.transport_stream_id", "1"),
					resource.TestCheckResourceAttr(resourceName, "multiplex_settings.0.maximum_video_buffer_delay_milliseconds", "1000"),
				),
			},
			{
				Config: testAccMultiplexConfig_update(rName, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMultiplexExists(ctx, resourceName, &multiplex),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrSet(resourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "multiplex_settings.0.transport_stream_bitrate", "1000001"),
					resource.TestCheckResourceAttr(resourceName, "multiplex_settings.0.transport_stream_reserved_bitrate", "1"),
					resource.TestCheckResourceAttr(resourceName, "multiplex_settings.0.transport_stream_id", "2"),
					resource.TestCheckResourceAttr(resourceName, "multiplex_settings.0.maximum_video_buffer_delay_milliseconds", "1000"),
				),
			},
		},
	})
}

func testAccMultiplex_updateTags(t *testing.T) {
	ctx := acctest.Context(t)
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	var multiplex medialive.DescribeMultiplexOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_medialive_multiplex.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPartitionHasService(names.MediaLiveEndpointID, t)
			testAccMultiplexesPreCheck(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.MediaLiveEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckMultiplexDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccMultiplexConfig_tags1(rName, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMultiplexExists(ctx, resourceName, &multiplex),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
				),
			},
			{
				Config: testAccMultiplexConfig_tags2(rName, "key1", "value1", "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMultiplexExists(ctx, resourceName, &multiplex),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccMultiplexConfig_tags1(rName, "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMultiplexExists(ctx, resourceName, &multiplex),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
		},
	})
}

func testAccMultiplex_disappears(t *testing.T) {
	ctx := acctest.Context(t)
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	var multiplex medialive.DescribeMultiplexOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_medialive_multiplex.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPartitionHasService(names.MediaLiveEndpointID, t)
			testAccMultiplexesPreCheck(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.MediaLiveEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckMultiplexDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccMultiplexConfig_basic(rName, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMultiplexExists(ctx, resourceName, &multiplex),
					acctest.CheckResourceDisappears(ctx, acctest.Provider, tfmedialive.ResourceMultiplex(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckMultiplexDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).MediaLiveClient()

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_medialive_multiplex" {
				continue
			}

			_, err := tfmedialive.FindInputSecurityGroupByID(ctx, conn, rs.Primary.ID)

			if tfresource.NotFound(err) {
				continue
			}

			if err != nil {
				return create.Error(names.MediaLive, create.ErrActionCheckingDestroyed, tfmedialive.ResNameMultiplex, rs.Primary.ID, err)
			}
		}

		return nil
	}
}

func testAccCheckMultiplexExists(ctx context.Context, name string, multiplex *medialive.DescribeMultiplexOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return create.Error(names.MediaLive, create.ErrActionCheckingExistence, tfmedialive.ResNameMultiplex, name, errors.New("not found"))
		}

		if rs.Primary.ID == "" {
			return create.Error(names.MediaLive, create.ErrActionCheckingExistence, tfmedialive.ResNameMultiplex, name, errors.New("not set"))
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).MediaLiveClient()

		resp, err := tfmedialive.FindMultiplexByID(ctx, conn, rs.Primary.ID)

		if err != nil {
			return create.Error(names.MediaLive, create.ErrActionCheckingExistence, tfmedialive.ResNameMultiplex, rs.Primary.ID, err)
		}

		*multiplex = *resp

		return nil
	}
}

func testAccMultiplexesPreCheck(ctx context.Context, t *testing.T) {
	conn := acctest.Provider.Meta().(*conns.AWSClient).MediaLiveClient()

	input := &medialive.ListMultiplexesInput{}
	_, err := conn.ListMultiplexes(ctx, input)

	if acctest.PreCheckSkipError(err) {
		t.Skipf("skipping acceptance testing: %s", err)
	}

	if err != nil {
		t.Fatalf("unexpected PreCheck error: %s", err)
	}
}

func testAccMultiplexBaseConfig() string {
	return `
data "aws_availability_zones" "available" {
  state         = "available"
  exclude_names = ["us-west-2-las-1a"] # not available in this zone

  filter {
    name   = "opt-in-status"
    values = ["opt-in-not-required"]
  }
}
`
}

func testAccMultiplexConfig_basic(rName string, start bool) string {
	return acctest.ConfigCompose(
		testAccMultiplexBaseConfig(),
		fmt.Sprintf(`
resource "aws_medialive_multiplex" "test" {
  name               = %[1]q
  availability_zones = [data.aws_availability_zones.available.names[0], data.aws_availability_zones.available.names[1]]

  multiplex_settings {
    transport_stream_bitrate                = 1000000
    transport_stream_id                     = 1
    transport_stream_reserved_bitrate       = 1
    maximum_video_buffer_delay_milliseconds = 1000
  }

  start_multiplex = %[2]t

  tags = {
    Name = %[1]q
  }
}
`, rName, start))
}

func testAccMultiplexConfig_update(rName string, start bool) string {
	return acctest.ConfigCompose(
		testAccMultiplexBaseConfig(),
		fmt.Sprintf(`
resource "aws_medialive_multiplex" "test" {
  name               = %[1]q
  availability_zones = [data.aws_availability_zones.available.names[0], data.aws_availability_zones.available.names[1]]

  multiplex_settings {
    transport_stream_bitrate                = 1000001
    transport_stream_id                     = 2
    transport_stream_reserved_bitrate       = 1
    maximum_video_buffer_delay_milliseconds = 1000
  }

  start_multiplex = %[2]t

  tags = {
    Name = %[1]q
  }
}
`, rName, start))
}

func testAccMultiplexConfig_tags1(rName, key1, value1 string) string {
	return acctest.ConfigCompose(
		testAccMultiplexBaseConfig(),
		fmt.Sprintf(`
resource "aws_medialive_multiplex" "test" {
  name               = %[1]q
  availability_zones = [data.aws_availability_zones.available.names[0], data.aws_availability_zones.available.names[1]]

  multiplex_settings {
    transport_stream_bitrate                = 1000000
    transport_stream_id                     = 1
    transport_stream_reserved_bitrate       = 1
    maximum_video_buffer_delay_milliseconds = 1000
  }

  tags = {
    %[2]q = %[3]q
  }
}
`, rName, key1, value1))
}

func testAccMultiplexConfig_tags2(rName, key1, value1, key2, value2 string) string {
	return acctest.ConfigCompose(
		testAccMultiplexBaseConfig(),
		fmt.Sprintf(`
resource "aws_medialive_multiplex" "test" {
  name               = %[1]q
  availability_zones = [data.aws_availability_zones.available.names[0], data.aws_availability_zones.available.names[1]]

  multiplex_settings {
    transport_stream_bitrate                = 1000000
    transport_stream_id                     = 1
    transport_stream_reserved_bitrate       = 1
    maximum_video_buffer_delay_milliseconds = 1000
  }

  tags = {
    %[2]q = %[3]q
    %[4]q = %[5]q
  }
}
`, rName, key1, value1, key2, value2))
}
