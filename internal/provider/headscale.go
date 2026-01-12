package main

import (
	"net/http"
	"time"
	"context"
	"fmt"
)

type Headscale struct {
	url string
	apikey string
}



func (h *Headscale) CheckAPI(ctx context.Context) (bool, time.Duration, error){
	
    

	req, err := http.NewRequestWithContext(ctx, "GET", h.url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return false,0, err
	}
	req.Header.Set("Authorization", "Bearer" + h.apikey)	

	



start := time.Now()
    

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Request error:", err)
		return false,0,err
	}
	defer resp.Body.Close()

	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Println("Response status:", resp.Status)
	
	return resp.StatusCode == 200, elapsed, nil


}	