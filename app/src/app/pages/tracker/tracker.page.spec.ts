import { ComponentFixture, TestBed } from '@angular/core/testing';
import { IonicModule } from '@ionic/angular';

import { TrackerPage } from './tracker.page';

describe('TrackerPage', () => {
  let component: TrackerPage;
  let fixture: ComponentFixture<TrackerPage>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [TrackerPage, IonicModule],
    }).compileComponents();

    fixture = TestBed.createComponent(TrackerPage);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
