import { Component, OnInit } from '@angular/core';
import { GraphqlService, Question } from '../graphql.service'
import {map} from 'rxjs/operators';
import { UserService } from '../user.service';

@Component({
  selector: 'app-questions',
  templateUrl: './questions.component.html',
  styleUrls: ['./questions.component.css']
})
export class QuestionsComponent implements OnInit {

  private questions: Array<Question>

  constructor(
    private graphqlService: GraphqlService,
    private userService: UserService
  ) { 
    this.questions = new Array<Question>();
  }

  ngOnInit() {
    let userID = this.userService.current().id
    this.graphqlService.queryQuestions(userID).subscribe( // get current user? global state management?
      (result) => { 
        this.questions = this.questions.concat(result.data.action.questions)
        console.log(this.questions)
      }
    )
  }

}
