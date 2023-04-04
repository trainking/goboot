package encrypt

import "testing"

func TestAes(t *testing.T) {
	tool := NewAesTool([]byte("66014775009e4106"), 16)
	if srcB, err := tool.Encrypt([]byte("Nihil in et et sea facilisis et erat duis aliquam lorem nulla kasd dolor kasd. In aliquyam sed ut ipsum dolor eos zzril ipsum. Vero suscipit stet gubergren eleifend eirmod sanctus ad quis. Iriure erat sit consetetur elitr et enim elitr vulputate assum rebum nulla erat feugiat nonumy tempor. Rebum labore iriure et in amet eu invidunt invidunt duis stet justo. Volutpat aliquyam aliquyam dolor amet clita duo ipsum. Ea lorem magna vero et gubergren dolor laoreet ut eirmod vero gubergren elitr amet justo. Dolor consequat et takimata velit praesent eirmod eos amet duo sit duis nostrud duo dolores no elit. Dolor at ut ut dolores clita nonumy te rebum at. Ut tincidunt veniam clita sed sadipscing et dolor. Lobortis et et. facer ipsum lorem. Dolore exerci at lobortis et ipsum duo erat blandit stet hendrerit gubergren. Lobortis sit erat takimata labore takimata assum ut ut gubergren diam vel invidunt erat takimata sed. Labore ipsum ut diam consequat ea dolores nisl lorem et eos at ea tempor amet. Sanctus vel quis gubergren stet ullamcorper no.")); err != nil {
		t.Error(err)
	} else {
		dB, err := tool.Decrypt(srcB)
		if err != nil {
			t.Error(err)
		}

		t.Logf(string(dB))
	}
}
