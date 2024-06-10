package payments

import "github.com/ShynaliOtep/paydal-api/models"

const Transaction_type_refill = "refill"
const Transaction_type_debit = "debit"

func TransactionBonus(userId uint, bonus float64, typeTr string) (err error) {

	/* пллучаем кошолек по userId если не
	найдем кошолек тогда создаем нового кошелка
	*/
	wallet, err := models.GetWalletByUserId(userId)
	if err != nil {
		wallet = models.Wallet{
			UserId:  userId,
			Balance: 0,
		}
		_, err = wallet.Create()
		if err != nil {
			return err
		}
	}

	transaction := models.Transaction{
		WalletId:    wallet.ID,
		Amount:      bonus,
		Type:        typeTr,
		BalanceType: "bonus",
	}
	_, err = transaction.Create()
	if err != nil {
		return err
	}

	err = topOfWalletBonus(wallet, bonus)
	if err != nil {
		return
	}
	return nil
}

func TransactionBalance(userId uint, balance float64, typeTr string) (err error) {

	/* пллучаем кошолек по userId если не
	найдем кошолек тогда создаем нового кошелка
	*/
	wallet, err := models.GetWalletByUserId(userId)
	if err != nil {
		wallet = models.Wallet{
			UserId:  userId,
			Balance: 0,
		}
		_, err = wallet.Create()
		if err != nil {
			return err
		}
	}

	transaction := models.Transaction{
		WalletId:    wallet.ID,
		Amount:      balance,
		Type:        typeTr,
		BalanceType: "money",
	}
	_, err = transaction.Create()
	if err != nil {
		return err
	}

	err = topOfWallet(wallet, balance)
	if err != nil {
		return
	}
	return nil
}
