# document_scanner.rb — Ruby версия

require 'opencv'
include OpenCV

def scan_document(input_path)
  image = CvMat.load(input_path)
  if image.nil?
    puts "❌ Не удалось загрузить"
    return
  end
  ratio = image.rows / 500.0
  resized = image.resize(CvSize.new((image.cols / ratio).to_i, 500))
  
  gray = resized.BGR2GRAY
  blurred = gray.GaussianBlur(Size.new(5,5), 0)
  edged = blurred.Canny(75, 200)

  contours = edged.find_contours
  sorted = contours.sort_by { |c| -c.area }.first(5)

  screen_cnt = nil
  sorted.each do |cnt|
    peri = cnt.arc_length(true)
    approx = cnt.approx_poly(0.02 * peri, true)
    if approx.size == 4
      screen_cnt = approx
      break
    end
  end
  if screen_cnt.nil?
    puts "❌ Контур не найден"
    return
  end
  # Коррекция пропущена
  out_path = input_path.sub(/\.[^.]+$/, '') + '_scanned.jpg'
  image.save(out_path)
  puts "💾 Сохранён: #{out_path}"
end

if ARGV.length < 1
  puts "Usage: ruby document_scanner.rb <image>"
  exit 1
end

input = ARGV[0]
puts "📄 Document Scanner (Ruby)"
puts "📂 Загружено: #{input}"
start = Time.now
scan_document(input)
puts "⏱️ Время: #{Time.now - start} сек."
