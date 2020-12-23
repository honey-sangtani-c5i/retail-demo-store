package main
import "fmt"
import "reflect"

func main() {
    fmt.Println("hello world")

    s := struct{ A int }{1} 
    field := reflect.ValueOf(s).Field(0) 
      
    // Use of IsZero() method 
    // fmt.Println(reflect.Value())
    fmt.Println(field.IsZero()) 

    var  b = 0
    
    fmt.Println(shouldIgnore(reflect.ValueOf(b)))
    

}

func shouldIgnore(v reflect.Value) bool {
	// if !ignoreEmpty {
	// 	return false
	// }
	
	// log.Printf(reflect.Value)
	// log.Printf(v)
	

	return v.IsZero()
}


