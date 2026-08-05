import type { ACTIVITY_TYPES } from '@/constants/user-activity'

export type BookmarkMetadata = {
  novela_title: string
  desc: string
}

export type UserNovelaLikeMetadata = {
  name: string
  desc: string
}

export type UserNovelaRatingMetadata = {
  name: string
  desc: string
  rating: number
}

export type AchievementMetadata = {
  name: string
  desc: string
}

export type RanobeAddMetadata = {
  ranobe_title: string
  author: string
}

export type CommentMetadata = {
  novela_title: string
  desc: string
}

interface BaseActivity {
  id: number
  target_id: number
  timestamp: string
}

interface BaseCreateActivity {
  user_id: number
  target_id: number
}

interface BookmarkActivity extends BaseActivity {
  activity_type: typeof ACTIVITY_TYPES.NOVELA_BOOKMARK
  metadata: BookmarkMetadata
}

interface RemoveBookmarkActivity extends BaseActivity {
  activity_type: typeof ACTIVITY_TYPES.NOVELA_BOOKMARK_REMOVE
  metadata: BookmarkMetadata
}

interface AchievementActivity extends BaseActivity {
  activity_type: typeof ACTIVITY_TYPES.ACHIEVEMENT_EARNED
  metadata: AchievementMetadata
}

interface RanobeActivity extends BaseActivity {
  activity_type: typeof ACTIVITY_TYPES.RANOBE_ADD
  metadata: RanobeAddMetadata
}

interface UserNovelaLikeActivity extends BaseActivity {
  activity_type: typeof ACTIVITY_TYPES.USER_NOVELA_LIKE
  metadata: UserNovelaLikeMetadata
}

interface UserNovelaRemoveLikeActivity extends BaseActivity {
  activity_type: typeof ACTIVITY_TYPES.USER_NOVELA_LIKE_REMOVE
  metadata: UserNovelaLikeMetadata
}

interface UserNovelaRatingActivity extends BaseActivity {
  activity_type: typeof ACTIVITY_TYPES.USER_NOVELA_RATING
  metadata: UserNovelaRatingMetadata
}

interface CommentActivity extends BaseActivity {
  activity_type: typeof ACTIVITY_TYPES.COMMENT_ADDED
  metadata: CommentMetadata
}

interface CommentRemoveActivity extends BaseActivity {
  activity_type: typeof ACTIVITY_TYPES.COMMENT_REMOVE
  metadata: CommentMetadata
}

interface CommentLikeActivity extends BaseActivity {
  activity_type: typeof ACTIVITY_TYPES.COMMENT_LIKE
  metadata: CommentMetadata
}

interface CommentUpdateActivity extends BaseActivity {
  activity_type: typeof ACTIVITY_TYPES.COMMENT_UPDATE
  metadata: CommentMetadata
}

interface CommentReportActivity extends BaseActivity {
  activity_type: typeof ACTIVITY_TYPES.COMMENT_REPORT
  metadata: CommentMetadata
}

// CREATE ACTIVITIES
interface CreateBookmarkActivity extends BaseCreateActivity {
  activity_type: typeof ACTIVITY_TYPES.NOVELA_BOOKMARK
  metadata: BookmarkMetadata
}

interface CreateRemoveBookmarkActivity extends BaseCreateActivity {
  activity_type: typeof ACTIVITY_TYPES.NOVELA_BOOKMARK_REMOVE
  metadata: BookmarkMetadata
}

interface CreateAchievementActivity extends BaseCreateActivity {
  activity_type: typeof ACTIVITY_TYPES.ACHIEVEMENT_EARNED
  metadata: AchievementMetadata
}

interface CreateRanobeActivity extends BaseCreateActivity {
  activity_type: typeof ACTIVITY_TYPES.RANOBE_ADD
  metadata: RanobeAddMetadata
}

interface CreateUserNovelaLike extends BaseCreateActivity {
  activity_type: typeof ACTIVITY_TYPES.USER_NOVELA_LIKE
  metadata: UserNovelaLikeMetadata
}

interface CreateUserNovelaRemoveLike extends BaseCreateActivity {
  activity_type: typeof ACTIVITY_TYPES.USER_NOVELA_LIKE_REMOVE
  metadata: UserNovelaLikeMetadata
}

interface CreateUserNovelaRating extends BaseCreateActivity {
  activity_type: typeof ACTIVITY_TYPES.USER_NOVELA_RATING
  metadata: UserNovelaRatingMetadata
}

interface CreateCommentActivity extends BaseCreateActivity {
  activity_type: typeof ACTIVITY_TYPES.COMMENT_ADDED
  metadata: CommentMetadata
}

interface CreateCommentRemoveActivity extends BaseCreateActivity {
  activity_type: typeof ACTIVITY_TYPES.COMMENT_REMOVE
  metadata: CommentMetadata
}

interface CreateCommentUpdateActivity extends BaseCreateActivity {
  activity_type: typeof ACTIVITY_TYPES.COMMENT_UPDATE
  metadata: CommentMetadata
}

interface CreateCommentLikeActivity extends BaseCreateActivity {
  activity_type: typeof ACTIVITY_TYPES.COMMENT_LIKE
  metadata: CommentMetadata
}

interface CreateCommentReportActivity extends BaseCreateActivity {
  activity_type: typeof ACTIVITY_TYPES.COMMENT_REPORT
  metadata: CommentMetadata
}

export type UserActivity =
  | BookmarkActivity
  | AchievementActivity
  | RanobeActivity
  | RemoveBookmarkActivity
  | UserNovelaLikeActivity
  | UserNovelaRemoveLikeActivity
  | UserNovelaRatingActivity
  | CommentActivity
  | CommentRemoveActivity
  | CommentLikeActivity
  | CommentUpdateActivity
  | CommentReportActivity

export type CreateUserActivity =
  | CreateBookmarkActivity
  | CreateAchievementActivity
  | CreateRanobeActivity
  | CreateRemoveBookmarkActivity
  | CreateUserNovelaLike
  | CreateUserNovelaRemoveLike
  | CreateUserNovelaRating
  | CreateCommentActivity
  | CreateCommentRemoveActivity
  | CreateCommentLikeActivity
  | CreateCommentUpdateActivity
  | CreateCommentReportActivity

export interface GetActivityResponse {
  message: UserActivity[]
}
