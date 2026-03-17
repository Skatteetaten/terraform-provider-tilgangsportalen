package provider

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestCreateNewEntraGroup(t *testing.T) {
	t.Parallel()

	// A timestamp is added to the name to avoid failure due to previous
	// test failures
	time := time.Now().Unix()
	name := fmt.Sprintf("[APPTEST] Test-Create_New_Entra_Group %d", time)
	newName := name + " new name"
	description := "APPTEST"

	var originalFieldValue string

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "tilgangsportalen_entra_group" "test" {
					name = "%s"
					description = "%s"
					inheritance_level = "User"
				}
				`, name, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tilgangsportalen_entra_group.test", "name", name),
					resource.TestCheckResourceAttr("tilgangsportalen_entra_group.test", "description", description),
					resource.TestCheckResourceAttr("tilgangsportalen_entra_group.test", "inheritance_level", "User"),
					resource.TestCheckResourceAttrSet("tilgangsportalen_entra_group.test", "object_id"),       // check that object_id is not an empty string
					resource.TestCheckResourceAttrSet("tilgangsportalen_entra_group.test", "entitlement_uid"), // check that entitlement_uid is not an empty string

					// Retrieve the object_id before name update
					func(s *terraform.State) error {
						rs := s.RootModule().Resources["tilgangsportalen_entra_group.test"]
						originalFieldValue = rs.Primary.Attributes["object_id"]
						return nil
					},
				),
			},
			// test import to state using ImportStateCheckFunc
			{
				ImportState:             true,
				ResourceName:            "tilgangsportalen_entra_group.test",
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"description", "inheritance_level", "entitlement_uid"},
			},
			// test update name
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "tilgangsportalen_entra_group" "test" {
					name = "%s"
					description = "%s"
					inheritance_level = "User"
				}
				`, newName, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tilgangsportalen_entra_group.test", "name", newName),
					resource.TestCheckResourceAttr("tilgangsportalen_entra_group.test", "description", description),
					resource.TestCheckResourceAttr("tilgangsportalen_entra_group.test", "inheritance_level", "User"),
					resource.TestCheckResourceAttrSet("tilgangsportalen_entra_group.test", "object_id"),       // check that object_id is not an empty string
					resource.TestCheckResourceAttrSet("tilgangsportalen_entra_group.test", "entitlement_uid"), // check that entitlement_uid is not an empty string
					// Verify that the group_id has NOT changed after name update
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources["tilgangsportalen_entra_group.test"]
						if !ok {
							return fmt.Errorf("resource not found in state")
						}
						newValue := rs.Primary.Attributes["object_id"]
						if newValue != originalFieldValue {
							return fmt.Errorf("object_id changed: was %s, now %s", originalFieldValue, newValue)
						}
						return nil
					},
				),
			},
		},
	})
}

func TestCreateNewEntraGroupWithSpecialCharacters(t *testing.T) {
	t.Parallel()

	// A timestamp is added to the name to avoid failure due to previous
	// test failures
	time := time.Now().Unix()
	name := fmt.Sprintf("[Group] Test-Create_New_Entra_Group %d", time)
	newName := name + " new name"
	//Testing special characters and line breaks in the group description.
	//The number of "\" characters changes in the expected description due to the use of GO´s Raw string literals in the input description.
	description := `<<-EOT
	Terraform_'acceptance'\n-_øåæ\tØÅÆ
	'!#$x%&/
	()[]{}'?!=(a){b}[c]@^*<>:,;.|
	\"test\" \t\"r\"ole
	EOT`
	expectedDescription := "Terraform_'acceptance'\\n-_øåæ\\tØÅÆ\n'!#$x%&/\n()[]{}'?!=(a){b}[c]@^*<>:,;.|\n\\\"test\\\" \\t\\\"r\\\"ole\n"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "tilgangsportalen_entra_group" "test" {
					name = "%s"
					description = %s
					inheritance_level = "User"
				}
				`, name, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tilgangsportalen_entra_group.test", "name", name),
					resource.TestCheckResourceAttr("tilgangsportalen_entra_group.test", "description", expectedDescription),
					resource.TestCheckResourceAttr("tilgangsportalen_entra_group.test", "inheritance_level", "User"),
				),
			},
			// test import to state using ImportStateCheckFunc
			{
				ImportState:             true,
				ResourceName:            "tilgangsportalen_entra_group.test",
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"description", "inheritance_level", "entitlement_uid"},
			},
			// test update name
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "tilgangsportalen_entra_group" "test" {
					name = "%s"
					description = %s
					inheritance_level = "User"
				}
				`, newName, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tilgangsportalen_entra_group.test", "name", newName),
					resource.TestCheckResourceAttr("tilgangsportalen_entra_group.test", "description", expectedDescription),
					resource.TestCheckResourceAttr("tilgangsportalen_entra_group.test", "inheritance_level", "User"),
				),
			},
		},
	})
}

func generateRandomHash(length int) (string, error) {
	// Generate random bytes
	bytes := make([]byte, length/2) // since hex encoding doubles the length
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	// Convert bytes to a hexadecimal string
	hash := hex.EncodeToString(bytes)

	// Truncate the hash to the desired length
	if len(hash) > length {
		hash = hash[:length]
	}

	return hash, nil
}

// test that group name length of 256 characters works
func TestCreateNewEntraGroupWithMaxNameLength(t *testing.T) {
	t.Parallel()

	description := "APPTEST"
	// generate a group name with 256 characters

	// create a group with a name of 256 characters
	maxNameLength := 256
	maxName, err := generateRandomHash(maxNameLength - 11)
	if err != nil {
		t.Fatalf("failed to generate random hash: %v", err)
	}
	maxName = fmt.Sprintf("[APPTEST] %s", maxName)

	newMaxName, err := generateRandomHash(maxNameLength - 8)
	if err != nil {
		t.Fatalf("failed to generate random hash: %v", err)
	}
	newMaxName = fmt.Sprintf("[Group] %s", newMaxName)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "tilgangsportalen_entra_group" "test" {
					name = "%s"
					description = "%s"
					inheritance_level = "User"
				}
				`, maxName, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tilgangsportalen_entra_group.test", "name", maxName),
					resource.TestCheckResourceAttr("tilgangsportalen_entra_group.test", "description", description),
					resource.TestCheckResourceAttr("tilgangsportalen_entra_group.test", "inheritance_level", "User"),
					resource.TestCheckResourceAttrSet("tilgangsportalen_entra_group.test", "object_id"),
					resource.TestCheckResourceAttrSet("tilgangsportalen_entra_group.test", "entitlement_uid"),
				),
			},
			// test import state
			{
				ImportState:             true,
				ResourceName:            "tilgangsportalen_entra_group.test",
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"description", "inheritance_level", "entitlement_uid"},
			},
			// test update name
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "tilgangsportalen_entra_group" "test" {
					name = "%s"
					description = "%s"
					inheritance_level = "User"
				}
				`, newMaxName, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tilgangsportalen_entra_group.test", "name", newMaxName),
					resource.TestCheckResourceAttr("tilgangsportalen_entra_group.test", "description", description),
					resource.TestCheckResourceAttr("tilgangsportalen_entra_group.test", "inheritance_level", "User"),
					resource.TestCheckResourceAttrSet("tilgangsportalen_entra_group.test", "object_id"),
					resource.TestCheckResourceAttrSet("tilgangsportalen_entra_group.test", "entitlement_uid"),
				),
			},
		},
	})
}

func TestCreateNewEntraGroupThatIsNotCreatedInEntra(t *testing.T) {
	t.Parallel()

	// A timestamp is added to the name to avoid failure due to previous
	// test failures
	time := time.Now().Unix()
	name := fmt.Sprintf("[TESTNOENTRA] Test-Create_New_Entra_Group_Not_In_Entra %d", time)
	description := "Group not created in Entra"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "tilgangsportalen_entra_group" "test" {
					name = "%s"
					description = "%s"
					inheritance_level = "User"
				}
				`, name, description),
				ExpectError: regexp.MustCompile("Unable to import entra group .*"),
			},
		},
	})
}
