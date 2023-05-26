import { Component, EnvironmentInjector, inject } from '@angular/core';
import { IonicModule } from '@ionic/angular';
import { CommonModule } from '@angular/common';
import { RouterOutlet } from '@angular/router';
import { TranslocoService } from '@ngneat/transloco';

@Component({
  selector: 'app-root',
  templateUrl: 'app.component.html',
  styleUrls: ['app.component.scss'],
  standalone: true,
  imports: [IonicModule, CommonModule, RouterOutlet],
})
export class AppComponent {
  public environmentInjector = inject(EnvironmentInjector);

  constructor(private transloco: TranslocoService) {
    const deviceLanguage = localStorage.getItem('language') || 'de';
    this.transloco.setDefaultLang(deviceLanguage);
    this.transloco.setActiveLang(deviceLanguage);
    localStorage.setItem('language', deviceLanguage);
  }
}
