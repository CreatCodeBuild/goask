import { Component, OnInit, Input } from '@angular/core';
import { Tag } from '../graphql.service';

@Component({
  selector: 'app-tags',
  templateUrl: './tags.component.html',
  styleUrls: ['./tags.component.css']
})
export class TagsComponent implements OnInit {

  @Input() private tags: Tag[]

  constructor(
    
  ) { }

  ngOnInit() {
  }

}
