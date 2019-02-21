import { Component, OnInit, Input } from '@angular/core';
import { Answer } from '../graphql.service';

@Component({
  selector: 'app-answer-detail',
  templateUrl: './answer-detail.component.html',
  styleUrls: ['./answer-detail.component.css']
})
export class AnswerDetailComponent implements OnInit {

  @Input() answer: Answer

  constructor() { }

  ngOnInit() {
  }

}
