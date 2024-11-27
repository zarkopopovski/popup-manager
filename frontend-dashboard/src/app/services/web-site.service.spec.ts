import { TestBed } from '@angular/core/testing';

import { WebSiteService } from './web-site.service';

describe('WebSiteService', () => {
  let service: WebSiteService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(WebSiteService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
