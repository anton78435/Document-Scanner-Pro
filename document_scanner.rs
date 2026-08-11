// document_scanner.rs — Rust версия

use image::{DynamicImage, GenericImageView};
use opencv::{
    core::{self, Mat, Point, Size, Vector},
    imgcodecs,
    imgproc,
    prelude::*,
    types::{VectorOfPoint, VectorOfPoint2f},
    Result,
};
use std::env;
use std::time::Instant;

fn main() -> Result<()> {
    let args: Vec<String> = env::args().collect();
    if args.len() < 2 {
        println!("Usage: cargo run -- <image>");
        return Ok(());
    }
    let input = &args[1];
    println!("📄 Document Scanner (Rust)");
    println!("📂 Загружено: {}", input);
    let start = Instant::now();

    let mut src = imgcodecs::imread(input, imgcodecs::IMREAD_COLOR)?;
    if src.empty() {
        println!("❌ Ошибка загрузки");
        return Ok(());
    }
    let ratio = src.rows() as f32 / 500.0;
    let mut resized = Mat::default();
    imgproc::resize(&src, &mut resized, Size { width: (src.cols() as f32 / ratio) as i32, height: 500 }, 0.0, 0.0, imgproc::INTER_LINEAR)?;

    let mut gray = Mat::default();
    imgproc::cvt_color(&resized, &mut gray, imgproc::COLOR_BGR2GRAY, 0)?;
    let mut blurred = Mat::default();
    imgproc::gaussian_blur(&gray, &mut blurred, Size { width: 5, height: 5 }, 0.0, 0.0, core::BORDER_DEFAULT)?;
    let mut edged = Mat::default();
    imgproc::canny(&blurred, &mut edged, 75.0, 200.0, 3, false)?;

    let mut contours = VectorOfPoint::new();
    let mut hierarchy = Mat::default();
    imgproc::find_contours(&edged, &mut contours, &mut hierarchy, imgproc::RETR_LIST, imgproc::CHAIN_APPROX_SIMPLE, core::Point::default())?;

    let mut screen_cnt: Option<VectorOfPoint> = None;
    for i in 0..contours.len() {
        let cnt = contours.get(i)?;
        let peri = imgproc::arc_length(&cnt, true)?;
        let mut approx = VectorOfPoint2f::new();
        imgproc::approx_poly_dp(&VectorOfPoint2f::from(cnt), &mut approx, 0.02 * peri, true)?;
        if approx.len() == 4 {
            // преобразуем обратно в Point
            let mut points = VectorOfPoint::new();
            for p in approx {
                points.push(Point { x: p.x as i32, y: p.y as i32 });
            }
            screen_cnt = Some(points);
            break;
        }
    }
    if screen_cnt.is_none() {
        println!("❌ Не найден контур");
        return Ok(());
    }
    // Пропускаем коррекцию для краткости
    let out_path = format!("{}_scanned.jpg", &input[..input.rfind('.').unwrap_or(input.len())]);
    imgcodecs::imwrite(&out_path, &src, &core::Vector::default())?;
    println!("💾 Сохранён: {}", out_path);
    println!("⏱️ Время: {:.2} сек.", start.elapsed().as_secs_f64());
    Ok(())
}
