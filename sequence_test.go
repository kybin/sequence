package sequence

import (
	"reflect"
	"testing"
)

func TestSplitter(t *testing.T) {
	cases := []struct {
		fname string
		want  []string
	}{
		{
			fname: "img.0001.exr",
			want:  []string{"img.", "0001", ".exr"},
		},
		{
			fname: "img_0001.exr",
			want:  []string{"img_", "0001", ".exr"},
		},
		{
			fname: "img0001.exr",
			want:  []string{"img", "0001", ".exr"},
		},
		{
			fname: "/a/b/c/img.0001.exr",
			want:  []string{"/a/b/c/img.", "0001", ".exr"},
		},
	}
	splitter := NewSplitter()
	for _, c := range cases {
		gotPre, gotDigits, gotPost, err := splitter.Split(c.fname)
		if err != nil {
			t.Fatalf("got err: %v", err)
		}
		got := []string{gotPre, gotDigits, gotPost}
		if !reflect.DeepEqual(got, c.want) {
			t.Fatalf("got: %q, want: %q", got, c.want)
		}
	}
}

func Test(t *testing.T) {
	cases := []struct {
		files []string
		want  string
	}{
		{
			files: []string{
				"/a/b/c/img.0001.exr",
				"/a/b/c/img.0002.exr",
				"/a/b/c/img.0003.exr",
				"/a/b/c/img.0004.exr",
				"/a/b/c/img.0098.exr",
				"/a/b/c/img.0099.exr",
				"/a/b/c/img.0100.exr",
				"/d/e/f/img.00001.exr",
			},
			want: "/a/b/c/img.####.exr 1-4 98-100\n/d/e/f/img.#####.exr 1",
		},
	}

	for _, c := range cases {
		man := NewManager(NewSplitter(), FmtSharp)
		for _, f := range c.files {
			err := man.Add(f)
			if err != nil {
				t.Fatalf("got error: %v", err)
			}
		}
		got := man.String()
		if got != c.want {
			t.Fatalf("got: %q, want: %q", got, c.want)
		}
	}
}
