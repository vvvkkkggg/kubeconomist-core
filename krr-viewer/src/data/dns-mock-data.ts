import type { DnsRecommendation } from '../types';

export const mockDnsRecommendations: DnsRecommendation[] = [
  {
    id: 'dnsa1qlkf6k5bad7ncrf',
    cloudId: 'b1g8vl7ekjq3a5d9f2m1',
    folderId: 'fo1e3k9lmn4p5q6r7s8',
    zoneId: 'dnsa1qlkf6k5bad7ncrf',
    isUsed: false,
  },
  {
    id: 'dnsb2qlkf6k5bad7ncrg',
    cloudId: 'b1g8vl7ekjq3a5d9f2m1',
    folderId: 'b1g8enqi42dal6t489ou',
    zoneId: 'dnsb2qlkf6k5bad7ncrg',
    isUsed: false,
  },
    {
    id: 'dnsc3qlkf6k5bad7ncrh',
    cloudId: 'b1g8vl7ekjq3a5d9f2m1',
    folderId: 'fo2a3b4c5d6e7f8g9h0',
    zoneId: 'dnsc3qlkf6k5bad7ncrh',
    isUsed: false,
  },
  // This one is in use and should be filtered out
  {
    id: 'dnsd4qlkf6k5bad7ncri',
    cloudId: 'b1g8vl7ekjq3a5d9f2m1',
    folderId: 'fo3c5d6e7f8g9h0i1j2',
    zoneId: 'dnsd4qlkf6k5bad7ncri',
    isUsed: true,
  },
]; 
