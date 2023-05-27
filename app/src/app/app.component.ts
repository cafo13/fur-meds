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
    this.setupLanguage();
    this.setupTheme();
  }

  setupLanguage() {
    const deviceLanguage = localStorage.getItem('language') || 'de';
    this.transloco.setDefaultLang(deviceLanguage);
    this.transloco.setActiveLang(deviceLanguage);
    localStorage.setItem('language', deviceLanguage);
  }

  setupTheme() {
    const appSettingsUseDarkTheme = localStorage.getItem('useDarkTheme');
    const deviceSettingsUseDarkTheme = window.matchMedia(
      '(prefers-color-scheme: dark)'
    ).matches;
    if (appSettingsUseDarkTheme !== null) {
      document.body.classList.toggle(
        'dark',
        JSON.parse(appSettingsUseDarkTheme)
      );
      localStorage.setItem('useDarkTheme', JSON.parse(appSettingsUseDarkTheme));
    } else {
      document.body.classList.toggle('dark', deviceSettingsUseDarkTheme);
      localStorage.setItem(
        'useDarkTheme',
        deviceSettingsUseDarkTheme.toString()
      );
    }
  }
}
