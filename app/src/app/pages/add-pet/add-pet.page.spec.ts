import { ComponentFixture, TestBed } from '@angular/core/testing';
import { IonicModule } from '@ionic/angular';

import { AddPetPage } from './add-pet.page';

describe('AddPetPage', () => {
  let component: AddPetPage;
  let fixture: ComponentFixture<AddPetPage>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [AddPetPage, IonicModule],
    }).compileComponents();

    fixture = TestBed.createComponent(AddPetPage);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
