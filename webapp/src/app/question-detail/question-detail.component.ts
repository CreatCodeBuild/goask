import { Component, OnInit } from '@angular/core';
import { GraphqlService, Question } from '../graphql.service';
import { ActivatedRoute } from '@angular/router';
import { UserService } from '../user.service';

@Component({
  selector: 'app-question-detail',
  templateUrl: './question-detail.component.html',
  styleUrls: ['./question-detail.component.css']
})
export class QuestionDetailComponent implements OnInit {

  private question: Question;

  constructor(
    private graphqlService: GraphqlService,
    private userService: UserService,
    private route: ActivatedRoute // todo? how to use it?
  ) { 

  }

  ngOnInit() {
    this.route.paramMap.subscribe(async (paramMap) => {
      let questionID = paramMap.get("id")
      let userID = this.userService.current().id
      let result = await this.graphqlService.queryQuestionDetail(userID, questionID).toPromise()
      this.question = result.data.action.question
    })  
  }

}
