package helper

import (
	"fmt"
	"log"
	"net/http"
)

func HandleError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	log.Printf("%+v\n", err)
	fmt.Fprintf(w, err.Error())
}
