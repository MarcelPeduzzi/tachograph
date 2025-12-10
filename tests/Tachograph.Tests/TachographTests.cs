using NUnit.Framework;
using System;

namespace Tachograph.Tests
{
    [TestFixture]
    public class TachographTests
    {
        [Test]
        public void Unmarshal_WithInsufficientData_ThrowsException()
        {
            // Arrange
            var data = new byte[] { 0x00 };

            // Act & Assert
            Assert.Throws<System.IO.InvalidDataException>(() => Tachograph.Unmarshal(data));
        }

        [Test]
        public void Unmarshal_WithNullData_ThrowsArgumentNullException()
        {
            // Act & Assert
            Assert.Throws<ArgumentNullException>(() => Tachograph.Unmarshal(null));
        }

        [Test]
        public void Parse_WithNullRawFile_ThrowsArgumentNullException()
        {
            // Act & Assert
            Assert.Throws<ArgumentNullException>(() => Tachograph.Parse(null));
        }

        [Test]
        public void Marshal_WithNullFile_ThrowsArgumentNullException()
        {
            // Act & Assert
            Assert.Throws<ArgumentNullException>(() => Tachograph.Marshal(null));
        }

        [Test]
        public void Anonymize_WithNullFile_ThrowsArgumentNullException()
        {
            // Act & Assert
            Assert.Throws<ArgumentNullException>(() => Tachograph.Anonymize(null));
        }
    }
}
