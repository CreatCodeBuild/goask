import { Component, OnInit } from '@angular/core';
import { UserService } from '../user.service';

@Component({
  selector: 'app-user-me',
  templateUrl: './user-me.component.html',
  styleUrls: ['./user-me.component.css']
})
export class UserMeComponent{

  constructor(
    private userService: UserService
  ) { }

  myID(): string {
    return this.userService.current().id
  }

}
