# Golang Gorm Transaction Manager

Simply implements DB nest transaction manager feature with gorm, Inspired by Laravel framework (Experimental)

# Using Example

```go
// bussiness logic code
func main(){

    // get *gorm.DB from your db manager
    db := core.NewDbManager().GetDb()

    // init transaction manger with the active db resource
    txm := core.NewTransactionManager(db)

    txm.Transaction(func() {

    	services.NewFinanceService().IncreaseBalance(txm, 18, 0.01)
    	services.NewFinanceService().IncreaseBalance(txm, 18, 0.04)
    	services.NewFinanceService().IncreaseBalance(txm, 18, 0.06)
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