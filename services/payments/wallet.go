package payments

import "github.com/ShynaliOtep/paydal-api/models"

func topOfWallet(wallet models.Wallet, sum float64) (err error) {
	wallet.Balance = wallet.Balance + sum
	_, err = wallet.Save()
	return err
}

func topOfWalletBonus(wallet models.Wallet, sum float64) (err error) {
	wallet.BonusBalance = wallet.BonusBalance + sum
	_, err = wallet.Save()
	return err
}
