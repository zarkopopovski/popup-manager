import { ComponentFixture, TestBed } from '@angular/core/testing';

import { DashboardWebsitesComponent } from './dashboard-websites.component';

describe('DashboardWebsitesComponent', () => {
  let component: DashboardWebsitesComponent;
  let fixture: ComponentFixture<DashboardWebsitesComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [DashboardWebsitesComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(DashboardWebsitesComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
