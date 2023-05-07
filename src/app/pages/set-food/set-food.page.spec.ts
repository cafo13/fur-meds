import { ComponentFixture, TestBed } from '@angular/core/testing';
import { IonicModule } from '@ionic/angular';

import { SetFoodPage } from './set-food.page';

describe('SetFoodPage', () => {
  let component: SetFoodPage;
  let fixture: ComponentFixture<SetFoodPage>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [SetFoodPage, IonicModule],
    }).compileComponents();

    fixture = TestBed.createComponent(SetFoodPage);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
