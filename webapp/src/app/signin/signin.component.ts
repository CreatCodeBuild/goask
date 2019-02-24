import { Component, OnInit } from '@angular/core';
import { NgForm } from '@angular/forms';
import { GraphqlService } from '../graphql.service';
import { UserService } from '../user.service';

@Component({
  selector: 'app-signin',
  templateUrl: './signin.component.html',
  styleUrls: ['./signin.component.css']
})
export class SigninComponent {

  constructor(
    private graphqlService: GraphqlService,
    private userService: UserService
  ) { }

  submit(form: NgForm) {
    console.log(form.value)
    this.graphqlService.querySignIn(
      this.userService.current().id, 
      form.value.userID, 
      form.value.userName
    ).subscribe((result)=>{

      if(result.data != null) {
        console.log(result.data.action.signIn)
        this.userService.User = result.data.action.signIn
      }

    })
  }

}
