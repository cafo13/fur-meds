import { Component } from '@angular/core';
import { IonicModule } from '@ionic/angular';
import { TranslocoModule } from '@ngneat/transloco';

@Component({
  selector: 'app-todo',
  templateUrl: 'todo.page.html',
  styleUrls: ['todo.page.scss'],
  standalone: true,
  imports: [IonicModule, TranslocoModule],
})
export class TodoPage {
  constructor() {}
}
