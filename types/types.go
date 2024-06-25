package types

type MnemonicBody struct {
	Mnemonic string `binding:"required" json:"mnemonic"`
}
