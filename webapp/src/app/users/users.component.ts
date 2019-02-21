import { Component, OnInit } from '@angular/core';
import { GraphqlService } from '../graphql.service';
import {map, filter} from 'rxjs/operators';

@Component({
  selector: 'app-users',
  templateUrl: './users.component.html',
  styleUrls: ['./users.component.css']
})
export class UsersComponent implements OnInit {
  users$: any;

  constructor(
    private userService: GraphqlService
  ) {
  }

  ngOnInit() {
    this.users$ = this.userService.queryUsers()
    .pipe(
      filter(res => res.loading === false),
      map(res => res.data.action.users)
    );
  }

}
