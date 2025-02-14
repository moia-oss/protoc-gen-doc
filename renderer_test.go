package gendoc_test

import (
	"os"
	"testing"

	. "github.com/moia-oss/protoc-gen-doc"
	"github.com/moia-oss/protokit"
	"github.com/moia-oss/protokit/utils"
	"github.com/stretchr/testify/require"
)

func TestRenderers(t *testing.T) {
	set, err := utils.LoadDescriptorSet("fixtures", "fileset.pb")
	require.NoError(t, err)

	tmp_dir := "./tmp"
	if os.Stat(tmp_dir); os.IsNotExist(err) {
		err = os.Mkdir(tmp_dir, os.ModePerm)
		require.NoError(t, err)
	}

	req := utils.CreateGenRequest(set, "Booking.proto", "Vehicle.proto")
	result := protokit.ParseCodeGenRequest(req)
	template := NewTemplate(result)

	for _, r := range []RenderType{
		RenderTypeDocBook,
		RenderTypeHTML,
		RenderTypeJSON,
		RenderTypeMarkdown,
	} {
		_, err := RenderTemplate(r, template, "")
		require.NoError(t, err)
	}
}

func TestNewRenderType(t *testing.T) {
	expected := []RenderType{
		RenderTypeDocBook,
		RenderTypeHTML,
		RenderTypeJSON,
		RenderTypeMarkdown,
	}

	supplied := []string{"docbook", "html", "json", "markdown"}

	for idx, input := range supplied {
		rt, err := NewRenderType(input)
		require.Nil(t, err)
		require.Equal(t, expected[idx], rt)
	}
}

func TestNewRenderTypeUnknown(t *testing.T) {
	rt, err := NewRenderType("/some/template.tmpl")
	require.Zero(t, rt)
	require.Error(t, err)
}
