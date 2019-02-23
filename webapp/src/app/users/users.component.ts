import { Component, OnInit } from '@angular/core';
import { GraphqlService, User } from '../graphql.service';

@Component({
  selector: 'app-users',
  templateUrl: './users.component.html',
  styleUrls: ['./users.component.css']
})
export class UsersComponent implements OnInit {
  users: User[];

  constructor(
    private userService: GraphqlService
  ) {
  }

  ngOnInit() {
    this.userService.queryUsers().subscribe((next)=>{
      this.users = next.data.action.users
    })
  }

}
