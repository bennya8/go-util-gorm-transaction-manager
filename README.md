# Golang Gorm Transaction Manager

Simply implements DB nest transaction manager feature with gorm, Inspired by Laravel framework (Experimental)

# Using Example

```go
// bussiness logic code example

func main(){

  // get *gorm.DB from your db manager
  db := core.NewDbManager().GetDb()

  // init transaction manger with the active db resource
  txm := core.NewTransactionManager(db)

  txm.Transaction(func() {

    result := services.NewFinanceService().IncreaseBalance(txm, 18, 0.01)
    // check service response state, if something goes wrong, just need to call panic() to trigger an error, all transactions will be rollback
    if !result.State {
      panic(result.Message)
    }
    
    result = services.NewFinanceService().IncreaseBalance(txm, 18, 0.04)
    if !result.State {
      panic(result.Message)
    }

    result = services.NewFinanceService().IncreaseBalance(txm, 18, 0.06)
    if !result.State {
      panic(result.Message)
    }
  
    // trigger an error, above transactions will automatically rollback
    panic("something happen")
  })
}

// service logic code
type FinanceService struct {}
func (f *FinanceService) IncreaseBalance(txm *core.TransactionManager, userId int, amount float32) (rs *base.ResultState) {

  db := txm.GetTx()

  defer func() {
    if r := recover(); r != nil {
      rs = base.NewResultState(false, fmt.Sprintf("%s", r), nil)
    }
  }()

  var finance entities.Finance

  err := db.Set("gorm:query_option", "FOR UPDATE").Where("user_id = ?", userId).First(&finance).Error
  if err != nil {
    panic(err)
  }
  if finance.ID <= 0 {
    panic("Fail, custom error message")
  }

  finance.Balance += amount

  err = db.Save(&finance).Error
  if err != nil {
    panic(err)
  }
  
  // nest call anthoer service 
  result = services.NewCreditService().IncreaseCredit(txm, 18, 0.04)
  if !result.State {
    panic(result.Message)
  }

  return base.NewResultState(true, "Success", nil)
}

type CreditService struct {}
func (f *CreditService) IncreaseCredit(txm *core.TransactionManager, userId int, amount float32) (rs *base.ResultState) {

  db := txm.GetTx()

  defer func() {
    if r := recover(); r != nil {
      rs = base.NewResultState(false, fmt.Sprintf("%s", r), nil)
    }
  }()

  var credit entities.Credit

  err := db.Set("gorm:query_option", "FOR UPDATE").Where("user_id = ?", userId).First(&credit).Error
  if err != nil {
    panic(err)
  }
  if credit.ID <= 0 {
    panic("Fail, custom error message")
  }

  credit.amount += amount

  err = db.Save(&credit).Error
  if err != nil {
    panic(err)
  }

  return base.NewResultState(true, "Success", nil)
}


// common response result pack for services
type ResultState struct {
  State   bool        `json:"state"`
  Message string      `json:"message"`
  Data    interface{} `json:"data"`
}

func NewResultState(state bool, message string, data interface{}) *ResultState {
  return &ResultState{
    State:   state,
    Message: message,
    Data:    data,
  }
}


```