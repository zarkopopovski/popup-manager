import { ComponentFixture, TestBed } from '@angular/core/testing';

import { DashboardPopupMessagesComponent } from './dashboard-popup-messages.component';

describe('DashboardPopupMessagesComponent', () => {
  let component: DashboardPopupMessagesComponent;
  let fixture: ComponentFixture<DashboardPopupMessagesComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [DashboardPopupMessagesComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(DashboardPopupMessagesComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
