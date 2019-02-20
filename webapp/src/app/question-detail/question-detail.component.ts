import { Component, OnInit } from '@angular/core';
import { GraphqlService, Question } from '../graphql.service';
import { ActivatedRoute } from '@angular/router';

@Component({
  selector: 'app-question-detail',
  templateUrl: './question-detail.component.html',
  styleUrls: ['./question-detail.component.css']
})
export class QuestionDetailComponent implements OnInit {

  private question: Question;

  constructor(
    private graphqlService: GraphqlService,
    private route: ActivatedRoute // todo? how to use it?
  ) { 

  }

  ngOnInit() {
    let self = this;
    this.graphqlService.queryQuestionDetail("1").subscribe(function next(result) {
      console.log(result)
      self.question = result.data.action.question
    })
  }

}
