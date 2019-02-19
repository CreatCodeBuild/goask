import { Component, OnInit } from '@angular/core';
import { UserService } from '../user.service';
import {map, filter} from 'rxjs/operators';
import { Observable } from 'apollo-link';

@Component({
  selector: 'app-users',
  templateUrl: './users.component.html',
  styleUrls: ['./users.component.css']
})
export class UsersComponent implements OnInit {
  users$: any;

  constructor(
    private userService: UserService
  ) {
  }

  ngOnInit() {
    this.users$ = this.userService.getAllUsers()
    .pipe(
      filter(res => res.loading === false),
      map(res => res.data.action.users)
    );
  }

}
