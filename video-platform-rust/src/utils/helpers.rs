use chrono::{DateTime, Utc};

pub fn format_duration(seconds: i32) -> String {
    let hours = seconds / 3600;
    let minutes = (seconds % 3600) / 60;
    let secs = seconds % 60;

    if hours > 0 {
        format!("{}:{:02}:{:02}", hours, minutes, secs)
    } else {
        format!("{}:{:02}", minutes, secs)
    }
}

pub fn format_number(num: i64) -> String {
    if num >= 1_000_000 {
        format!("{:.1}M", num as f64 / 1_000_000.0)
    } else if num >= 1_000 {
        format!("{:.1}K", num as f64 / 1_000.0)
    } else {
        num.to_string()
    }
}

pub fn time_ago(dt: DateTime<Utc>) -> String {
    let now = Utc::now();
    let diff = now.signed_duration_since(dt);

    if diff.num_seconds() < 60 {
        "just now".to_string()
    } else if diff.num_minutes() < 60 {
        format!("{} minutes ago", diff.num_minutes())
    } else if diff.num_hours() < 24 {
        format!("{} hours ago", diff.num_hours())
    } else if diff.num_days() < 30 {
        format!("{} days ago", diff.num_days())
    } else if diff.num_days() < 365 {
        format!("{} months ago", diff.num_days() / 30)
    } else {
        format!("{} years ago", diff.num_days() / 365)
    }
}
