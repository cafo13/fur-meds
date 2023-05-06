import { ComponentFixture, TestBed } from '@angular/core/testing';
import { IonicModule } from '@ionic/angular';

import { PetPage } from './pet.page';

describe('PetPage', () => {
  let component: PetPage;
  let fixture: ComponentFixture<PetPage>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [PetPage, IonicModule],
    }).compileComponents();

    fixture = TestBed.createComponent(PetPage);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
