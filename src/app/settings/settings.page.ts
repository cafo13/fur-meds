import { Component } from '@angular/core';
import { IonicModule } from '@ionic/angular';

@Component({
  selector: 'app-settings',
  templateUrl: 'settings.page.html',
  styleUrls: ['settings.page.scss'],
  standalone: true,
  imports: [IonicModule],
})
export class SettingsPage {
  usingSystemDarkTheme: boolean;

  constructor() {
    this.usingSystemDarkTheme = window.matchMedia(
      '(prefers-color-scheme: dark)'
    ).matches;
  }

  toggleDarkTheme(event: any): void {
    // document.body.classList.toggle('dark', event.detail.checked);
    // this.usingSystemDarkTheme = event.detail.checked
    console.log('to be implented: toggle dark theme here');
  }
}
