import type { VpcRecommendation } from "../types";

export const mockVpcRecommendations: VpcRecommendation[] = [
  {
    id: '1',
    cloudId: 'b1g8vl7ekjq3a5d9f2m1',
    folderId: 'b1g8enqi42dal6t489ou',
    ipAddress: '84.252.142.128',
    isUsed: false,
    isReserved: true,
  },
  {
    id: '2',
    cloudId: 'b1g8vl7ekjq3a5d9f2m1',
    folderId: 'fo2a3b4c5d6e7f8g9h0',
    ipAddress: '194.87.237.59',
    isUsed: false,
    isReserved: true,
  },
  {
    id: '3',
    cloudId: 'b1g8vl7ekjq3a5d9f2m1',
    folderId: 'fo3c5d6e7f8g9h0i1j2',
    ipAddress: '5.255.255.77',
    isUsed: false,
    isReserved: true,
  },
  {
    id: '4',
    cloudId: 'b1g8vl7ekjq3a5d9f2m1',
    folderId: 'fo3c5d6e7f8g9h0i1j2',
    ipAddress: '95.163.231.64',
    isUsed: false,
    isReserved: true,
  },
  {
    id: '5',
    cloudId: 'b1g8vl7ekjq3a5d9f2m1',
    folderId: 'fo4d7e8f9g0h1i2j3k4',
    ipAddress: '213.180.204.62',
    isUsed: false,
    isReserved: true,
  },
  // This one is used, so it should be filtered out in the view
  {
    id: '6',
    cloudId: 'b1g8vl7ekjq3a5d9f2m1',
    folderId: 'b1g8enqi42dal6t489ou',
    ipAddress: '77.88.55.88',
    isUsed: true,
    isReserved: true,
  },
];
