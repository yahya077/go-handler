# 🚧 gohandler
This is simplified request handler package. Main purpose of this package isolating services from controller. Best feature of this package is you can handle request separately from the service part.

## 💈 Full Example

```go
// request.go
type Request struct {
	gohandler.BaseRequest
}

// overrided method
func (r Request) GetSchema() interface{} {
	return r.RequestSchema.(*CustomRequest)
}

type CustomRequest struct {
	Username string  `json:"username" validate:"required"`
}

func NewRequest() gohandler.IRequest {
	return Request{
		gohandler.BaseRequest{
			RequestSchema: new(CustomRequest),
		},
	}
}

```

```go
// controller.go
type MyCustomHandler struct {
}

func (receiver MyCustomHandler) LocalBinding(ctx *fiber.Ctx) error {
	fmt.Println("you can bind parameters to context")
	return nil
}

func Get(ctx *fiber.Ctx) error {
	h := gohandler.New(ctx)
	h.WithHandler(MyCustomHandler{})   //optional
	h.WithRequest(myCustomRequest.New()) //optional

  // it will call your method, will bind ctx as parameters
	return h.Handle(myIsolatedService.New(), "Get") 
}

```

### 🧭 Quick Note

You can bind your struct to Fiber's [Locals](https://docs.gofiber.io/api/ctx#locals) on LocalBinding method. For example;
```go
func (h Handler) LocalBinding(ctx *fiber.Ctx) error {
	id := ctx.Params("id") // get id from param
  
	myStruct, _ := repository.New().FindByID(id) // find row from db
  
	ctx.Locals("myStruct", &myStruct) // bind your data to locals
  
	return nil
}
```

## 📫&nbsp; Have a question? Want to chat? Ran into a problem?

#### *Website [yahyahindioglu.com](https://yahyahindioglu.com)*

#### *[LinkedIn](https://www.linkedin.com/in/yahyahindioglu/) Account*

## License
This project is available under the [MIT](https://opensource.org/licenses/mit-license.php) license.
