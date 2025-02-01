package epicmodule

import "ZakkBob/AskDave/backend/epicmodule"

func TestImportedNum(t *testing.T) {
	if ImportedNum != 2 {
		t.ErrorF("Errored LOL")
	}
}