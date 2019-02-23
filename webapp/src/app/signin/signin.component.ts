import { Component, OnInit } from '@angular/core';
import { NgForm } from '@angular/forms';
import { GraphqlService } from '../graphql.service';

@Component({
  selector: 'app-signin',
  templateUrl: './signin.component.html',
  styleUrls: ['./signin.component.css']
})
export class SigninComponent {

  constructor(
    private graphqlService: GraphqlService
  ) { }

  submit(form: NgForm) {
    console.log(form.value)
    // this.graphqlService.querySignIn(form.value.userID, form.value.userName)
  }

}
