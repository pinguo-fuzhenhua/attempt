package main

import "fmt"

func main() {
	// html := fmt.Sprintf(`<div>  <div>   尊敬的，您好！  </div>  <div style="padding: 8px 40px 8px 50px;">   <p>您于提交的邮箱验证，本次验证码为 ，为了保证账号安全，验证码有效期为5分钟。请确认为本⼈操作，切勿向他⼈泄露，感谢您的理解与使⽤。</p >  </div>  <div>   <p>此邮箱为系统邮箱，请勿回复。</p >  </div>  </div>`)
	// msg := gomail.NewMessage()
	// msg.SetAddressHeader("From", "fuzhenhua@camera360.com", "XXX")
	// msg.SetHeader("To", "fuzhenhua@camera360.com")
	// msg.SetHeader("Subject", "用户反馈——【用户姓名】")
	// msg.SetBody("text/html", html)

	// d := gomail.NewDialer("smtp.exmail.qq.com", 25, "fuzhenhua@camera360.com", "Y2iaciej")
	// if err := d.DialAndSend(msg); err != nil {
	// 	fmt.Println(err.Error())
	// }
	a := struct{}{}
	fmt.Printf("%p %v %T", &a, a, a)
	fmt.Printf("c %T %v %p", c, c, &c)
}
