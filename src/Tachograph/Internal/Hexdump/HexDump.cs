using System;
using System.Text;

namespace Tachograph.Internal.Hexdump
{
    /// <summary>
    /// Hexdump utilities for debugging binary data.
    /// </summary>
    internal static class HexDump
    {
        /// <summary>
        /// Creates a hexadecimal dump of binary data.
        /// </summary>
        /// <param name="data">The data to dump</param>
        /// <param name="offset">Starting offset in the data</param>
        /// <param name="length">Number of bytes to dump (0 for all)</param>
        /// <returns>Hexadecimal string representation</returns>
        internal static string Dump(byte[] data, int offset = 0, int length = 0)
        {
            if (data == null) return string.Empty;
            
            int actualLength = length == 0 ? data.Length - offset : Math.Min(length, data.Length - offset);
            var sb = new StringBuilder();
            
            for (int i = 0; i < actualLength; i += 16)
            {
                sb.Append($"{offset + i:X8}  ");
                
                // Hex values
                for (int j = 0; j < 16; j++)
                {
                    if (i + j < actualLength)
                        sb.Append($"{data[offset + i + j]:X2} ");
                    else
                        sb.Append("   ");
                    
                    if (j == 7) sb.Append(" ");
                }
                
                sb.Append(" |");
                
                // ASCII representation
                for (int j = 0; j < 16 && i + j < actualLength; j++)
                {
                    byte b = data[offset + i + j];
                    sb.Append(b >= 32 && b < 127 ? (char)b : '.');
                }
                
                sb.AppendLine("|");
            }
            
            return sb.ToString();
        }
    }
}
