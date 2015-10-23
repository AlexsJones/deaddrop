package image

import (
  "image"
  "image/draw"
  "image/png"
  "image/jpeg"
  "io/ioutil"
  "log"
  "os"
  "github.com/AlexsJones/deaddrop/Godeps/_workspace/src/github.com/golang/freetype"
)
func getFormat(file *os.File) string {
  bytes := make([]byte, 4)
  n, _ := file.ReadAt(bytes, 0)
  if n < 4 { return "" }
  if bytes[0] == 0x89 && bytes[1] == 0x50 && bytes[2] == 0x4E && bytes[3] == 0x47 { 
    return "png" }
    if bytes[0] == 0xFF && bytes[1] == 0xD8 { return "jpg" }
    if bytes[0] == 0x47 && bytes[1] == 0x49 && bytes[2] == 0x46 && bytes[3] == 0x38 { 
      return "gif" }
      if bytes[0] == 0x42 && bytes[1] == 0x4D { return "bmp" }
      return ""
    }
  
    func GenerateWaterMark(inputImage string, watermarkText string) (string,error) {

      /* target image path */
      targetimg, _ := os.Open(inputImage)

      var targetImage image.Image

      switch getFormat(targetimg) {

      case "png":
        targetImage, _ = png.Decode(targetimg)

      case "jpg":
        targetImage, _ = jpeg.Decode(targetimg)
      default:
        log.Println("Format not excepted")
        os.Exit(1)
      }

      defer targetimg.Close()
      // Read the font data.
      fontBytes, err := ioutil.ReadFile("luxisr.ttf")
      if err != nil {
        log.Println(err)
        return "",err
      }
      f, err := freetype.ParseFont(fontBytes)
      if err != nil {
        log.Println(err)
        return "", err
      }

      // Initialize the context.
      fg, _ := image.Black, image.White

      rgba := image.NewRGBA(image.Rect(0, 0, 
      targetImage.Bounds().Size().X, targetImage.Bounds().Size().Y))

      draw.Draw(rgba, rgba.Bounds(), image.Transparent, image.ZP, draw.Src)

      c := freetype.NewContext()
      c.SetDPI(72)
      c.SetFont(f)
      c.SetFontSize(12)
      c.SetClip(rgba.Bounds())
      c.SetDst(rgba)
      c.SetSrc(fg)


      current_offset_y := 0
      current_offset_x := 0

      //totalWidth := targetImage.Bounds().Size().X
      totalHeight := targetImage.Bounds().Size().Y

      // intervalX := 30
      intervalY := 20

      for current_offset_y < totalHeight {
        // Draw the text.
        pt := freetype.Pt(current_offset_x, current_offset_y+int(c.PointToFixed(12)>>6))
        _, err = c.DrawString(watermarkText, pt)
        if err != nil {
          log.Println(err)
          return "", err
        }
        current_offset_y += intervalY
      }

      log.Println("Generated in memory watermark")

      ib := targetImage.Bounds()
      m := image.NewRGBA(ib)

      offset := image.Pt(0,0)

      draw.Draw(m, ib, targetImage, image.ZP, draw.Src)
      draw.Draw(m, rgba.Bounds().Add(offset),rgba, image.ZP, draw.Over)

      finalWaterMarkedImage,_ := os.Create("watermarked_" + inputImage)
      jpeg.Encode(finalWaterMarkedImage,m, &jpeg.Options{jpeg.DefaultQuality})
      defer finalWaterMarkedImage.Close()

      return "watermarked_" + inputImage,nil
}
