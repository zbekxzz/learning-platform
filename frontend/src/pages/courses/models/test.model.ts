export interface TestItem {
  id: number;
  module_id: number;
}

export type QuestionType = 'mcq' | 'open' | 'matching' | string;

export interface TestPair {
  left: string;
}

export interface TestQuestion {
  id: number;
  type: QuestionType;
  question_text: string;
  // For matching type:
  pairs?: TestPair[];
  options?: string[];
}

export interface TestAnswerOption {
  id: number;
  text: string;
}

export interface TestResponse {
  test: TestItem;
  questions: TestQuestion[];
  answers: Record<number, TestAnswerOption[]>;
}

export interface TestSubmitPayload {
  test_id: number;
  answers: Record<number, any>; // maps question_id to response (number for mcq, string for open, Record<string, string> for matching)
}

export interface TestSubmitResponse {
  score: number;
  passed?: boolean; // We might get a passed boolean or we determine success by HTTP 200
  feedback?: string;
}
