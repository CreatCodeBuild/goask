import { Component, OnInit, Input } from '@angular/core';
import { GraphqlService, User } from '../graphql.service';
import { UserService } from '../user.service';

@Component({
  selector: 'app-user-detail',
  templateUrl: './user-detail.component.html',
  styleUrls: ['./user-detail.component.css']
})
export class UserDetailComponent implements OnInit {

  private user: User

  @Input() private userID: string

  constructor(
    private graphqlService: GraphqlService,
    private userService: UserService
  ) { }

  ngOnInit() {
    console.log(this.userService.current().id, this.userID)
    let obs = this.graphqlService.queryUser(this.userService.current().id, this.userID)
    obs.subscribe(
      (value) => {
        console.log(value)
        this.user = value.data.action.user
      },
      (error) => {
        console.log(error)
      },
    )

  }

}
