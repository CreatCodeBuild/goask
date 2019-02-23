import { Component, OnInit, Input } from '@angular/core';
import { GraphqlService, User } from '../graphql.service';
import { UserService } from '../user.service';

@Component({
  selector: 'app-user-summary',
  templateUrl: './user-summary.component.html',
  styleUrls: ['./user-summary.component.css']
})
export class UserSummaryComponent implements OnInit {

  @Input() user: User

  constructor(
    private graphqlService: GraphqlService,
    private userService: UserService
  ) { }

  ngOnInit() {}

}
